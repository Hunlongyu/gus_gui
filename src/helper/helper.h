#ifndef HELPER_H
#define HELPER_H
#pragma once
#include <windows.h>

bool CreateDeviceD3D(HWND hWnd);
void CleanupDeviceD3D();
void CreateRenderTarget();
void CleanupRenderTarget();

#endif
