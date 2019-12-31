set PULL_WORKSPACE=%1%
set PULL_BRANCH=%2%
set DIST_PATH=%3%
set RUN_NAME=%4%

set cmdpath=%CD%

echo "Pulling Source Code"
cd %PULL_WORKSPACE%\%PULL_BRANCH%
git fetch -q
git checkout -q --force origin/%PULL_BRANCH%

go mod download
go build  -ldflags "-s -w" -o bin\main

mkdir %DIST_PATH%
copy bin\main %DIST_PATH%\%RUN_NAME%
echo "Finished"

cd %cmdpath%