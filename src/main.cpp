#include <iostream>

#include "imgui_impl_win32.h"
#include "imgui_impl_dx11.h"
#include "globals/globals.h"
#include "helper/helper.h"
#include "resource/resource.h"
#include "views/home/home.h"
#include "windows/win.h"

int APIENTRY WinMain(HINSTANCE hCurrentInstance, HINSTANCE hPreviousInstance, LPSTR lpcmdline, int nCmdShow)
{
    HWND hwnd = CreateApplicationWindow();

    if (!CreateDeviceD3D(hwnd))
    {
        CleanupDeviceD3D();
        UnregisterWindowClass();
        return 1;
    }

    // Show the window
    ::ShowWindow(hwnd, SW_SHOWDEFAULT);
    ::UpdateWindow(hwnd);

    initResource();
    homeInit(hwnd);

    {
        /*char email[] = "MyEmailAddress@example.com ";
        auto lala = GetEmailSha(email);
        ImGui::Text(lala);
        delete[] lala;*/
    }

    {
        /*const char* sha = "205e460b479e2e5b48aec07710c08d50";
        char* imgUrl = GetImageUrl(sha);
        std::cout << "Image URL: " << imgUrl << "\n";
        delete[] imgUrl;*/
    }

    // Main loop
    bool done = false;
    while (!done)
    {
        MSG msg;
        while (::PeekMessage(&msg, nullptr, 0U, 0U, PM_REMOVE))
        {
            TranslateMessage(&msg);
            ::DispatchMessage(&msg);
            if (msg.message == WM_QUIT)
            {
                done = true;
            }
        }
        if (done)
        {
            break;
        }

        // Handle window resize (we don't resize directly in the WM_SIZE handler)
        if (g_ResizeWidth != 0 && g_ResizeHeight != 0)
        {
            CleanupRenderTarget();
            g_pSwapChain->ResizeBuffers(0, g_ResizeWidth, g_ResizeHeight, DXGI_FORMAT_UNKNOWN, 0);
            g_ResizeWidth = g_ResizeHeight = 0;
            CreateRenderTarget();
        }

        // Start the Dear ImGui frame
        ImGui_ImplDX11_NewFrame();
        ImGui_ImplWin32_NewFrame();
        ImGui::NewFrame();

        // views code here
        if (show_home)
        {
            homeUI();
        }
        if (show_manage)
        {
            //
        }
        if (show_profile)
        {
            // 
        }
        if (show_settings)
        {
            //
        }

        if (show_demo)
        {
            ImGui::ShowDemoWindow(&show_demo);
        }

        // Rendering
        ImGui::Render();
        g_pd3dDeviceContext->OMSetRenderTargets(1, &g_mainRenderTargetView, nullptr);
        constexpr float clearColor[4] = { 0.00f, 0.00f, 0.00f, 1.00f }; // RGBA
        g_pd3dDeviceContext->ClearRenderTargetView(g_mainRenderTargetView, clearColor);
        ImGui_ImplDX11_RenderDrawData(ImGui::GetDrawData());
        g_pSwapChain->Present(1, 0); // Present with vsync
    }

    // Cleanup
    homeDestroy();
    CleanupDeviceD3D();
    UnregisterWindowClass();
    return 0;
}
