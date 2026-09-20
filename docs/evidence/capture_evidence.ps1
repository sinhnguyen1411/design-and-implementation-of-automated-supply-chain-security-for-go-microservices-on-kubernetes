Write-Host "=================================================" -ForegroundColor Cyan
Write-Host "KICH BAN 1: HARNESS BENCHMARK" -ForegroundColor Yellow
Write-Host "=================================================" -ForegroundColor Cyan
python scripts/test_dependency_firewall_harness.py
Write-Host ""
Write-Host ">>> THUC NGHIEM 1 DA XONG!" -ForegroundColor Green
Write-Host ">>> ANH HAY DUNG WIN+SHIFT+S DE CHUP LAI MAN HINH TERMINAL NAY." -ForegroundColor Yellow
Write-Host ">>> Luu anh vao: docs/evidence/powershell_harness.png" -ForegroundColor Yellow
Pause

Clear-Host
Write-Host "=================================================" -ForegroundColor Cyan
Write-Host "KICH BAN 2: CHAN TRUC TIEP GOI DOC HAI (LIVE BLOCKING)" -ForegroundColor Yellow
Write-Host "=================================================" -ForegroundColor Cyan
cd services/user-service
curl.exe -i http://localhost:8080/api/v1/dependency-proxy/4a882f42-bdb4-49a0-b5fe-0509eb92faa0/go/github.com/boltdb-go/bolt/@v/v1.3.1.info
cd ../..
Write-Host ""
Write-Host ">>> THUC NGHIEM 2 DA XONG (Du kien bao loi 403 Forbidden do DevGuard chan)!" -ForegroundColor Green
Write-Host ">>> ANH HAY DUNG WIN+SHIFT+S DE CHUP LAI MAN HINH TERMINAL NAY." -ForegroundColor Yellow
Write-Host ">>> Luu anh vao: docs/evidence/powershell_blocking.png" -ForegroundColor Yellow
Pause

Clear-Host
Write-Host "=================================================" -ForegroundColor Cyan
Write-Host "KICH BAN 3: SO SANH 3 MOI TRUONG (A/B/C)" -ForegroundColor Yellow
Write-Host "=================================================" -ForegroundColor Cyan
python scripts/run_phase3_comparative_evaluation.py
Write-Host ""
Write-Host ">>> THUC NGHIEM 3 DA XONG!" -ForegroundColor Green
Write-Host ">>> ANH HAY DUNG WIN+SHIFT+S DE CHUP LAI MAN HINH TERMINAL NAY." -ForegroundColor Yellow
Write-Host ">>> Luu anh vao: docs/evidence/powershell_phase3.png" -ForegroundColor Yellow
Pause

Write-Host "=================================================" -ForegroundColor Cyan
Write-Host "HOAN TAT! ANH HAY NHAN LAI CHO EM TREN KHUNG CHAT DE EM CAP NHAT BAO CAO NHE!" -ForegroundColor Green
