param (
    [switch]$Help
)

if ($Help) {
    Write-Host "Usage: .\test_ai_scs_local.ps1"
    Write-Host "Mô tả: Script tự động triển khai và kiểm thử kiến trúc Enterprise AI Supply Chain."
    exit
}

Write-Host "=== BẮT ĐẦU KIỂM THỬ AI SUPPLY CHAIN SECURITY (INIT-CONTAINER) ===" -ForegroundColor Cyan

# 1. Kiểm tra kết nối Kubernetes
Write-Host "`n[1/5] Kiểm tra kết nối Kubernetes cluster..." -ForegroundColor Yellow
$k8sStatus = kubectl cluster-info 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "LỖI: Không thể kết nối tới Kubernetes. Vui lòng bật Docker Desktop hoặc Kind." -ForegroundColor Red
    exit 1
}
Write-Host "OK: Kubernetes đang chạy." -ForegroundColor Green

# 2. Tạo Namespace
Write-Host "`n[2/5] Đảm bảo namespace 'stock-trading' tồn tại..." -ForegroundColor Yellow
kubectl create namespace stock-trading --dry-run=client -o yaml | kubectl apply -f -

# 3. Triển khai AI Backend (qwen-llm-service)
Write-Host "`n[3/5] Áp dụng các cấu hình Kubernetes cho qwen-llm-service..." -ForegroundColor Yellow
kubectl apply -f deploy/kubernetes/qwen-llm-service/
Write-Host "Đã apply Deployment, Service và HPA." -ForegroundColor Green

# 4. Giám sát Init-Container (Zero-Trust Fetching)
Write-Host "`n[4/5] Chờ Init-Container chạy và xác minh chữ ký mô hình (Cosign)..." -ForegroundColor Yellow
Write-Host "Bạn có thể mở một Terminal khác và gõ lệnh sau để xem log thực tế:" -ForegroundColor Gray
Write-Host "kubectl logs -f deployment/qwen-llm-service -c init-model-verifier -n stock-trading" -ForegroundColor Gray

# Chờ một chút để pod được schedule
Start-Sleep -Seconds 5

# Đợi Pod ở trạng thái Ready
Write-Host "Đang đợi Main Container (Ollama) khởi động thành công (Có thể mất 1-2 phút)..."
kubectl wait --for=condition=ready pod -l app=qwen-llm-service -n stock-trading --timeout=120s

if ($LASTEXITCODE -eq 0) {
    Write-Host "OK: Qwen LLM Service đã sẵn sàng!" -ForegroundColor Green
} else {
    Write-Host "CẢNH BÁO: Pod chưa Ready sau 120s. Vui lòng kiểm tra lại tài nguyên (RAM/CPU)." -ForegroundColor Red
    exit 1
}

# 5. Kiểm thử Business Logic (risk-service)
Write-Host "`n[5/5] Port-forward và chạy Unit Test của risk-service..." -ForegroundColor Yellow
Write-Host "Mở port-forward ngầm (background job)..."
Start-Job -Name "PortForwardOllama" -ScriptBlock {
    kubectl port-forward svc/qwen-llm-service -n stock-trading 11434:11434
} | Out-Null

Start-Sleep -Seconds 3 # Đợi port-forward ổn định

Write-Host "Chạy go test trong services/risk-service..."
Push-Location services/risk-service
go test ./internal/risk -v
$testResult = $LASTEXITCODE
Pop-Location

# Dọn dẹp
Write-Host "Đóng port-forward..."
Stop-Job -Name "PortForwardOllama"
Remove-Job -Name "PortForwardOllama"

if ($testResult -eq 0) {
    Write-Host "`n=== KIỂM THỬ THÀNH CÔNG TỐT ĐẸP! ===" -ForegroundColor Green
    Write-Host "Kiến trúc Init-Container đã bảo vệ và tải mô hình xuất sắc." -ForegroundColor Green
} else {
    Write-Host "`n=== KIỂM THỬ THẤT BẠI ===" -ForegroundColor Red
    Write-Host "Vui lòng xem lại log báo lỗi bên trên." -ForegroundColor Red
}
