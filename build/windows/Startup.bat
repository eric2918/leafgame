@echo build and run

@echo call build
@call:build gateway .\bin\conf\gateway.json
@goto:eof

:build
@set name=%~1
@set config=%~2
@cd .\cmd\%name%\
@go build -o .\bin\%name%.exe .\cmd\%name%\main.go
@start .\bin\%name%.exe config
@goto:eof
