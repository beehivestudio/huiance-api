package crypt

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"time"

	"github.com/astaxie/beego/logs"
)

/******************************************************************************
 **函数名称: Md5Sum
 **功    能: MD5加密处理
 **输入参数:
 **     s: 被加密处理的字串
 **输出参数: NONE
 **返    回: 加密字串
 **实现描述:
 **注意事项:
 **作    者: # Qifeng.zou # 2018.04.26 15:35:08 #
 ******************************************************************************/
func Md5Sum(str string) string {
	h := md5.New()

	h.Write([]byte(str))

	return string(hex.EncodeToString(h.Sum(nil)))
}

////////////////////////////////////////////////////////////////////////////////

type Md5EncodeCtx struct {
	h      hash.Hash
	cipher string
}

func CreateEncodeCtx(passwd string) *Md5EncodeCtx {
	encodeCtx := &Md5EncodeCtx{}
	encodeCtx.h = md5.New()
	encodeCtx.cipher = passwd
	return encodeCtx
}

func md5CipherEncode(ctx *Md5EncodeCtx, sourceText string) string {
	ctx.h.Write([]byte(ctx.cipher))
	cipherHash := fmt.Sprintf("%x", ctx.h.Sum(nil))
	ctx.h.Reset()
	inputData := []byte(sourceText)
	loopCount := len(inputData)
	outData := make([]byte, loopCount)
	for i := 0; i < loopCount; i++ {
		outData[i] = inputData[i] ^ cipherHash[i%32]
	}
	return string(outData)
}

func Md5Encode(ctx *Md5EncodeCtx, sourceText string) string {
	ctx.h.Write([]byte(time.Now().Format("2015-12-22 15:04:05")))
	noise := fmt.Sprintf("%x", ctx.h.Sum(nil))
	ctx.h.Reset()
	inputData := []byte(sourceText)
	loopCount := len(inputData)
	outData := make([]byte, loopCount*2)

	for i, j := 0, 0; i < loopCount; i, j = i+1, j+1 {
		outData[j] = noise[i%32]
		j++
		outData[j] = inputData[i] ^ noise[i%32]
	}
	return base64.StdEncoding.EncodeToString([]byte(md5CipherEncode(ctx, fmt.Sprintf("%s", outData))))
}

func Md5Decode(ctx *Md5EncodeCtx, sourceText string) string {
	buf, err := base64.StdEncoding.DecodeString(sourceText)
	if err != nil {
		fmt.Println("Decode(%q) failed: %v", sourceText, err)
		return ""
	}
	inputData := []byte(md5CipherEncode(ctx, fmt.Sprintf("%s", buf)))
	loopCount := len(inputData)
	outData := make([]byte, loopCount)
	for i, j := 0, 0; i < loopCount; i, j = i+2, j+1 {
		outData[j] = inputData[i] ^ inputData[i+1]
	}
	return string(outData)
}

/******************************************************************************
 **函数名称: Md5ByFilePath
 **功    能: 根据文件地址，对文件MD5加密处理
 **输入参数:
 **     filePath: 被加密处理的文件地址
 **输出参数: NONE
 **返    回:
 **     md5Str: 加密字符串
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-5-25 15:15:42 #
 ******************************************************************************/
func Md5ByFilePath(filePath string) (md5Str string, err error) {

	logs.Info("Md5ByFilePath request param. filePath: %s", filePath)

	md5 := md5.New()
	file, err := os.Open(filePath)
	if err != nil {
		logs.Error("Md5ByFilePath open file failed! err: %s", err.Error())
		return "", err
	}
	defer file.Close()

	written, err := io.Copy(md5, file)
	if err != nil {
		logs.Error("Md5ByFilePath copy file failed! err: %s", err.Error())
		return "", err
	}
	logs.Info("Md5ByFilePath written: %d", written)

	md5Str = hex.EncodeToString(md5.Sum(nil))
	logs.Info("Md5ByFilePath response param. "+filePath+"的md5为: %s", md5Str)

	return md5Str, nil
}
