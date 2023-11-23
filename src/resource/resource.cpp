#include "resource.h"
#include "imgui.h"
#include "../helper/helper.h"

ID3D11ShaderResourceView* img_menu = nullptr;
ID3D11ShaderResourceView* img_setting = nullptr;

void initResource()
{
    const bool menu_ret = LoadTextureFromFile("assets/images/menu.png", &img_menu);
    IM_ASSERT(menu_ret);

    const bool setting_ret = LoadTextureFromFile("assets/images/setting.png", &img_setting);
    IM_ASSERT(setting_ret);
}
