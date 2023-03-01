@echo off
setlocal

SET GIT_MASK_WRAPPER=%~f0
SET GIT_MASK_WRAPPER_DIR=%~dp0

FOR /f "delims=" %%G IN ('where git') DO (
    IF /I NOT ["%%~dpG"] == ["%GIT_MASK_WRAPPER_DIR%"] (
        SET GIT_MASK_REAL_GIT_PROGRAM=%%~fG
    )
)

CALL "%GIT_MASK_REAL_GIT_PROGRAM%" mask run %*
