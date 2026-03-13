# HDGP 合规测试脚本（需先在一终端运行: go run ./cmd/hdgp-engine）
# 用法: .\scripts\run-conftest.ps1  或在项目根目录: powershell -File scripts\run-conftest.ps1
Set-Location $PSScriptRoot\..
go run ./cmd/hdgp-conftest
exit $LASTEXITCODE
