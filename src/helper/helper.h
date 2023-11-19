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

bool LoadTextureFromFile(const char* filename, ID3D11ShaderResourceView** out_srv, int* out_width, int* out_height);
bool LoadTextureFromFile(const char* filename, ID3D11ShaderResourceView** out_srv);

#endif
