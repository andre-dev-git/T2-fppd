@echo off
echo Compilando servidor RPC...
go build -o server.exe main.go
if %errorlevel% neq 0 (
    echo Erro na compilacao do servidor
    exit /b 1
)

echo Compilacao concluida com sucesso!
echo.
echo Para executar o servidor:
echo   server.exe [porta]
echo   Exemplo: server.exe 8080
echo.
echo Servidor RPC pronto para receber conexoes de clientes!
echo.
pause
