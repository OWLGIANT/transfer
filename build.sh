EXECUTE_NAME="bpy_transfer"
echo "########################## build ${EXECUTE_NAME} start ##########################"
CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -o $EXECUTE_NAME -tags=prd
upx -9 $EXECUTE_NAME
echo "########################## build ${EXECUTE_NAME} end ##########################"
sleep 5  # 休眠 5 秒