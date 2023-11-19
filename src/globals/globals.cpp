#include "globals.h"

ID3D11Device* g_pd3dDevice = nullptr;
ID3D11DeviceContext* g_pd3dDeviceContext = nullptr;
IDXGISwapChain* g_pSwapChain = nullptr;
UINT g_ResizeWidth = 0, g_ResizeHeight = 0;
ID3D11RenderTargetView* g_mainRenderTargetView = nullptr;
HWND hwnd = nullptr;


bool show_home{true};
bool show_manage{false};
bool show_profile{false};
bool show_settings{false};

bool show_demo{true};
bool show_hello{false};
