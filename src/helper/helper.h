#ifndef HELPER_H
#define HELPER_H
#pragma once
#include <windows.h>
#define STB_IMAGE_IMPLEMENTATION
#include <d3d11.h>

bool CreateDeviceD3D(HWND hWnd);
void CleanupDeviceD3D();
void CreateRenderTarget();
void CleanupRenderTarget();

// 从本地文件加载图片
bool LoadTextureFromFile(const char* filename, ID3D11ShaderResourceView** out_srv, int* out_width, int* out_height);
bool LoadTextureFromFile(const char* filename, ID3D11ShaderResourceView** out_srv);

// 获取 email 的 sha256 准备通过接口获取头像地址
void TrimWhitespace(char* str);
void ToLower(char* str);
char* SHA256Str(const char* str);
char* GetEmailSha(char* email);

// 通过接口获取头像地址
char* GetImageUrl(const char* sha);
#endif
