// jquery selectors

jQuery.fn.reverse = function() {
    return this.pushStack(this.get().reverse(), arguments);
};

treeSelector = '#tree';

menuSelector = '#off-canvas-menu';
menuToggleSelector = '#sidebar-toggle';


var mmenu = {};
var jqtree = {};
var nodes = [];
var initialMenuPath = [];

var defaultTitle = "(undefined)";

var scrollbarOptions = {
    theme: "minimal",
    scrollbarPosition: "inside",
	scrollInertia: 0
};



// data functions
function getTree() {
    return (typeof tocData === "undefined" ? [] : tocData);
}

function getLanguageData() {
    return (typeof languagesData === "undefined" ? [] : languagesData);
}

function getMenuData() {
    return (typeof menuData === "undefined" ? [] : menuData);
}


function isTouchDevice() {
    //return Modernizr.touch;
    return (('ontouchstart' in window) ||
        (navigator.maxTouchPoints > 0) ||
        (navigator.msMaxTouchPoints > 0));
}


function setActiveLanguageMenu( $item ) {
    
            var lang = $item.attr("data-value");
            var curLang = getCurrentLanguage();

            if(lang == curLang) {
                    $item.addClass("checked") 
            }
}

// helper function for language menu
function switchLanguage(newLang) {

    var curLang = getCurrentLanguage();

    if (curLang != newLang) {

        var curPath = current_path = window.location.pathname.split('/').pop();

        var newLocation = '../' + newLang + '/' + curPath;

        window.location.href = newLocation;

    }
}


// gets the current language 
function getCurrentLanguage() {

    return $('body').attr('data-culture');

}

function getCurrentNodeId() {

    return $("body").attr("data-node-id");

}

function getCurrentMenuPanel() {

    var currentNodeId = getCurrentNodeId();
    var currentLevel = nodes[currentNodeId];

    if (currentLevel && currentLevel.parent) {
        return $("#panel_" + currentLevel.parent.node.id);

    } else {
        return $("#panel_root");
    }

}

function getDummyTreeNode() {

    return {
        label: "loading...",
        id: "loading",
        href: undefined,
        children: []
    };

}


function initializeNodesDictionary(treeNodes, parent) {

    treeNodes.forEach(function (node) {

        // store node to nodes dictionary
        nodes[node.id] = {
            node: node,
            parent: parent
        };

        if (node.nodes && node.nodes.length > 0) {
            initializeNodesDictionary(node.nodes, nodes[node.id]);
        }
    })
}



function appendMenuLevel(parent, parentId, nodes) {


    if (nodes && nodes.length) {

        var list = $("<ul id='panel_" + parentId + "'></ul>");
        list.appendTo(parent);


        if (parentId == "root") {

			var menuHeadingValue = typeof menuHeading === "undefined" ? "(Menu)" : menuHeading;
            var menuHeadingElement = $("<li subheading='1'><a>" + menuHeadingValue + "</a></li>");
            menuHeadingElement.appendTo(list);

            // add navigation triplet home, index, glossary
            menuData.forEach(function (menuItem) {

                var listItem = $("<li entry-id='entry_" + menuItem.id + "' " + (getCurrentNodeId() === menuItem.id ? "current='1'" : "") + "><a href='" + menuItem.href + "'>" + (menuItem.name.length ? menuItem.name : defaultTitle) + "</a></li>");

                listItem.appendTo(list)

            });

            // append languages menu
            if (getLanguageData().length > 1) {

                var languagesMenuItem = $("<li><a>" + languagesHeading + "</a></li>");
                languagesMenuItem.appendTo(list);


                var languagesMenu = $("<ul id='panel_Languages" + parentId + "'></ul>");
                languagesMenu.appendTo(languagesMenuItem);

                getLanguageData().forEach(function (languageDataItem) {

                    var languagesMenuEntry = $("<li aspect='" + languageDataItem.data + "' " + (getCurrentLanguage() == languageDataItem.data ? "checked='1'" : "") + "><a>" + languageDataItem.name + "</a></li>");
                    if (getCurrentLanguage() != languageDataItem.data) {
                        languagesMenuEntry.bind('click', function (e) {
                            var lang = $(e.currentTarget).attr("aspect");
                            switchLanguage(lang);
                        });
                    }

                    languagesMenuEntry.appendTo(languagesMenu);

                });

            }

			var tocHeadingValue = typeof tocHeading === "undefined" ? "(Table of contents)" : tocHeading;
            var tocHeadingElement = $("<li subheading='1'><a>" + tocHeadingValue + "</a></li>");
            tocHeadingElement.appendTo(list);



        }


        nodes.forEach(function (node) {

            var listItem = $("<li entry-id='entry_" + node.id + "' " + (getCurrentNodeId() === node.id ? "current='1'" : "") + "><a href='" + node.href + "'>" + (node.text.length ? node.text : defaultTitle) + "</a></li>");

            listItem.appendTo(list);

            if (node.nodes && node.nodes.length > 0) {

                if (!isOnInitialPath(node.id)) {

                    var childList = $("<ul id='panel_" + node.id + "' panel-id='panel_" + node.id + "'></ul>");
                    childList.appendTo(listItem);

                    var defaultItem = $("<li class='loading'><a>loading...</a></li>");
                    defaultItem.appendTo(childList);

                } else {

                    appendMenuLevel(listItem, node.id, node.nodes);

                }

            }

        });
    }

}

function loadMenuPanel($panel) {


    var nodeId = $panel.attr("id").replace("panel_", "");
    var node = nodes[nodeId].node;

    var list = $('#' + $panel.attr("id") + " > ul")

    if (list.find("li.loading").length > 0) {

        list.find("li.loading").remove();

        node.nodes.forEach(function (child) {

            var listItem = $("<li entry-id='entry_" + child.id + "' " + (getCurrentNodeId() === child.id ? "current='1'" : "current='0'") + "><a href='" + child.href + "'>" + (child.text.length ? child.text : defaultTitle) + "</a></li>");

            listItem.appendTo(list);

            if (child.nodes && child.nodes.length > 0) {

                var childList = $("<ul id='panel_" + child.id + "' panel-id='panel_" + child.id + "'></ul>");
                childList.appendTo(listItem);

                var defaultItem = $("<li class='loading'><a>loading...</a></li>");
                defaultItem.appendTo(childList);

                mmenu.init(childList);

            }

        });

        mmenu.update();

    }
}


function calculateInitialSubTree(level) {


    if (level && level.node && level.node.id) {
        initialMenuPath.push(level.node.id);
    }

    if (level && level.parent) {
        calculateInitialSubTree(level.parent);
    }

}

function isOnInitialPath(nodeId) {

    var foundNodesOnPath = initialMenuPath.filter(function (pathId) {
        return pathId == nodeId;
    });

    return (foundNodesOnPath.length > 0);

}


function initializeMenu() {

    var menu = $(menuSelector);

    var menuOptions = {
        "slidingSubmenus": isTouchDevice(),
        "navbar": {
            "add": true,
            "title": projectTitle,
            "titleLink": "parent"
        }
    }

    menu.addClass("loading");

    setTimeout(function () {

        var currentNodeId = getCurrentNodeId();
        var currentMenuLevel = nodes[currentNodeId];

        calculateInitialSubTree(currentMenuLevel);

        // append tree data to menu
        appendMenuLevel(menu, 'root', getTree());

        $(menu).mmenu(menuOptions);

        mmenu = $(menu).data("mmenu");

        mmenu.openPanel(getCurrentMenuPanel());
        mmenu.setSelected($("entry_" + getCurrentNodeId()));

        mmenu.bind("openPanel", function ($panel) {

            loadMenuPanel($panel);

        });

        setTimeout(function () {

            menu.removeClass("loading");
        }, 500);

    }, 1);


}



function openInitialSubTree(currentId) {


    var treeContainer = $('#tree-container');

    var currentNode = nodes[currentId];

    if (currentNode && currentNode.parent) {

        initialMenuPath.forEach(function (nodeId) {

            var node = treeContainer.tree('getNodeById', nodeId);
            treeContainer.tree('openNode', node);

        });
    }

    var currentNode = treeContainer.tree('getNodeById', getCurrentNodeId());
    treeContainer.tree('selectNode', currentNode);


}

function intializeSubTree(parent, nodes) {

    nodes.forEach(function (node) {

        var treeNode = {
            label: node.text,
            id: node.id,
            href: node.href,
            children: []
        };

        if (node.nodes && node.nodes.length) {

            if (isOnInitialPath(node.id)) {

                 intializeSubTree(treeNode.children, node.nodes);
            } else {
                treeNode.children.push(getDummyTreeNode());
            }

        }

        parent.push(treeNode);

    });

}


function insertNavMenu() {
    
    var additionalClass = "visible-sm";
    var menuDropdown = $(".navbar-dropdown");
	var menuHeadingValue = typeof menuHeading === "undefined" ? "(Menu)" : menuHeading;
    var menuHeader = $("<li class='dropdown-header " + additionalClass + "'>" + menuHeadingValue + "</li>");
    if( menuDropdown ) {
    
        $(".navbar-menu li").each( function() {

           var item = $( this ).clone();
            
            item.addClass( additionalClass );
            item.removeClass( "active" );
            

            
            item.prependTo( menuDropdown );


        });
        
        menuHeader.prependTo(menuDropdown);
    }
}

function initializeTree() {


    var firstTreeLevel = [];

    var currentNodeId = getCurrentNodeId();
    var currentMenuLevel = nodes[currentNodeId];

    calculateInitialSubTree(currentMenuLevel);

     getTree().forEach(function (node) {

        var treeNode = {
            label: node.text,
            id: node.id,
            href: node.href,
            children: []
        };

        if (node.nodes && node.nodes.length) {

            if (isOnInitialPath(node.id)) {
                intializeSubTree(treeNode.children, node.nodes);
            } else {
                treeNode.children.push(getDummyTreeNode());
            }

        }

        firstTreeLevel.push(treeNode);

    });

    var menu = $(menuSelector);

	var tocHeadingValue = typeof tocHeading === "undefined" ? "(Table of contents)" : tocHeading;
    var treeTitle = $("<div id='tree-title'>" + tocHeadingValue + "</div>").appendTo(menu);
    var treeContainer = $("<div id='tree-container'></div>").appendTo(menu);

    var options = {

        data: firstTreeLevel,
        slide: false,
        closedIcon: '&#xe602',
        openedIcon: '&#xe600',
        onCreateLi: function (node, $li) {

            $li.find('.jqtree-title').attr("title", node.name);
            
            if( node.class ) {
                $li.addClass( node.class );
            }
        }
    }

    var jqtree = treeContainer.tree(options);

    openInitialSubTree();

    var scrollbar = treeContainer.mCustomScrollbar(scrollbarOptions);

    jqtree.bind(
        'tree.click',
        function (event) {

            var node = event.node;

            if (node.href) {

                window.location = node.href;
                event.preventDefault();

            }

        }
    );

    jqtree.bind(
        'tree.open',
        function (event) {

            var node = event.node;

            var treeNode = event.node ; // jqtree.tree('getNodeById', node.id); 

 
            var realNode = nodes[node.id].node;

            var data = {
                label: realNode.text,
                id: realNode.id,
                href: realNode.href
            };

            var children = realNode.nodes.map(function (realChildNode) {

                var childNode = {
                    label: realChildNode.text,
                    id: realChildNode.id,
                    href: realChildNode.href,
                    children: []
                };

                if (realChildNode.nodes && realChildNode.nodes.length) {
                    childNode.children.push(getDummyTreeNode());
                }

                return childNode;


            });

            jqtree.tree('loadData', children, treeNode);



        }
    );


}

$(function () {


    $("body").addClass("loading");

    // get device class
    var device = (isTouchDevice() ? "mobile" : "desktop");

    $("body").addClass(device);
    
    setTimeout(function () {
        $("#off-canvas-menu").addClass("animated");
        $("#content-container").addClass("animated");
    }, 1000);
    
    

    $("." + device + "-hidden-sm").addClass("hidden-sm");
    $("." + device + "-hidden-md").addClass("hidden-md");
    $("." + device + "-hidden-xs").addClass("hidden-xs");
    $("." + device + "-hidden-lg").addClass("hidden-lg");

    $("." + device + "-visible-sm").addClass("visible-sm");
    $("." + device + "-visible-md").addClass("visible-md");
    $("." + device + "-visible-xs").addClass("visible-xs");
    $("." + device + "-visible-lg").addClass("visible-lg");

    // initialize nodes dictionary
    initializeNodesDictionary(getTree(), undefined);
    initializeNodesDictionary(getMenuData(), undefined);

    if (isTouchDevice()) {

        // intialize multi level push menu for touch devices
        initializeMenu();

    } else {

        // intialize tree for pointer devices
        initializeTree();
        insertNavMenu();

    }

    // initialize navbar menu toggle 
    $(menuToggleSelector).click(function () {

        var container = $(".off-canvas");
        var activeClass = "showMenu";

        if (!container.hasClass(activeClass)) {
            container.addClass(activeClass);
        } else {
            container.removeClass(activeClass);
        };

    });


    // set current language => main menu
    $('.navbar').find('.langMenuItem').each(function() {
       setActiveLanguageMenu( $( this ) );
    });
    
    
    // activate popupovers for lexicon links
    var popoverOptions = {
        container: "body",
        placement: "bottom",
        trigger: "hover",
        delay: {
            show: 200,
            hide: 100
        }
    };

    $('.lexicon-link').popover(popoverOptions);


    $("body").removeClass("loading");


});

$(window).resize(function () {


});