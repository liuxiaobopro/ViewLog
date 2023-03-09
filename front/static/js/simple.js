layui.use(['tree', 'code', 'dropdown'], function () {
    layui.code()
    var tree = layui.tree,
        layer = layui.layer,
        $ = layui.jquery,
        form = layui.form,
        util = layui.util

    //#region 获取隐藏域
    var activeSshId = $("#activeSshId").val(); // 当前激活的sshId
    //#endregion

    //#region 监听运行模式
    form.on('radio(mode)', function (data) {
        var mode = data.value;
        if (mode == "1") {
            $(".do_search").hide();
        } else {
            $(".do_search").show();
        }
    });
    //#endregion

    //#region 监听ssh多选框
    form.on('checkbox(ssh)', function (data) {
        let loading = layer.load(3)

        sshId = data.value;
        checked = data.elem.checked;
        isActive = checked ? 2 : 1

        $.ajax({
            async: true,
            // url: '/api/ssh',
            url: '/api/ssh/active',
            type: 'PUT',
            data: JSON.stringify(
                {
                    "id": parseInt(sshId),
                    "isActive": isActive
                }
            ),
            dataType: 'json',
            timeout: 30000,
            success: successCallback,
            error: errorCallback,
            complete: completeCallback
        })

        function successCallback(res) {
            console.log("res:", res);
            if (res.code == 0) {
                location.reload();
            } else {
                layer.msg(res.msg, { icon: 2, time: 2000 });
            }
        }

        function errorCallback(err, status) {
            console.error("err:", err)
            layer.msg(err.responseJSON.msg, { icon: 2, time: 2000 });
        }

        function completeCallback(xhr, status) {
            console.log('Ajax请求已结束。');
            layer.close(loading)
        }
    });
    //#endregion

    //#region 添加SSH
    $("#add_ssh").click(function () {
        layer.open({
            type: 2,
            title: '添加SSH',
            shadeClose: true,
            shade: 0.8,
            offset: '20%',
            area: ['30%', '45%'],
            content: 'ssh_add'
        });
    })
    //#endregion

    //#region 添加文件夹
    $("#add_folder").click(function () {
        layer.open({
            type: 2,
            title: '添加文件夹',
            shadeClose: true,
            shade: 0.8,
            offset: '20%',
            area: ['30%', '45%'],
            content: 'folder_add'
        });
    })
    //#endregion

    //#region ssh item鼠标悬浮
    $(".ssh-item1 .ssh_item").hover(function () {
        $(this).next().show();
    })
    $(".ssh-item1 .icons").hover(function () {
        $(this).show();
    })
    $(".ssh-item1 .ssh_item").mouseleave(function () {
        $(this).next().hide();
    })
    $(".ssh-item1 .icons").mouseleave(function () {
        $(this).hide();
    })
    //#endregion

    //#region 删除ssh
    $(".ssh-item1 .del_ssh").click(function () {
        var sshId = $(this).attr("data-id");
        var sshName = $(this).attr("data-name");
        layer.confirm('正在删除【' + sshName + '】,确定删除该SSH？', {
            icon: 3,
            btn: ['确定', '取消'] //按钮
        }, function () {
            $.ajax({
                async: true,
                url: '/api/ssh',
                type: 'DELETE',
                data: JSON.stringify(
                    {
                        "id": parseInt(sshId)
                    }
                ),
                dataType: 'json',
                timeout: 30000,
                success: successCallback,
                error: errorCallback,
                complete: completeCallback
            })

            function successCallback(res) {
                console.log("res:", res);
                if (res.code == 0) {
                    layer.msg(res.msg, { icon: 1, time: 1000 }, function () {
                        location.reload();
                    });
                } else {
                    layer.msg(res.msg, { icon: 2, time: 2000 });
                }
            }

            function errorCallback(err, status) {
                console.error("err:", err)
            }

            function completeCallback(xhr, status) {
                console.log('Ajax请求已结束。');
            }
        }, function () {
            // layer.msg('已取消', { icon: 1, time: 1000 });
        });
    })
    //#endregion

    //#region 展示文件夹
    data = [{
        title: '一级1'
        , id: 1
        , field: 'name1'
        , checked: true
        , spread: true
        , children: [{
            title: '二级1-1 可允许跳转'
            , id: 3
            , field: 'name11'
            , href: 'https://www.layui.com/'
            , children: [{
                title: '三级1-1-3'
                , id: 23
                , field: ''
                , children: [{
                    title: '四级1-1-3-1'
                    , id: 24
                    , field: ''
                    , children: [{
                        title: '五级1-1-3-1-1'
                        , id: 30
                        , field: ''
                    }, {
                        title: '五级1-1-3-1-2'
                        , id: 31
                        , field: ''
                    }]
                }]
            }, {
                title: '三级1-1-1'
                , id: 7
                , field: ''
                , children: [{
                    title: '四级1-1-1-1 可允许跳转'
                    , id: 15
                    , field: ''
                    , href: 'https://www.layui.com/doc/'
                }]
            }, {
                title: '三级1-1-2'
                , id: 8
                , field: ''
                , children: [{
                    title: '四级1-1-2-1'
                    , id: 32
                    , field: ''
                }]
            }]
        }, {
            title: '二级1-2'
            , id: 4
            , spread: true
            , children: [{
                title: '三级1-2-1'
                , id: 9
                , field: ''
                , disabled: true
            }, {
                title: '三级1-2-2'
                , id: 10
                , field: ''
            }]
        }, {
            title: '二级1-3'
            , id: 20
            , field: ''
            , children: [{
                title: '三级1-3-1'
                , id: 21
                , field: ''
            }, {
                title: '三级1-3-2'
                , id: 22
                , field: ''
            }]
        }]
    }, {
        title: '一级2'
        , id: 2
        , field: ''
        , spread: true
        , children: [{
            title: '二级2-1'
            , id: 5
            , field: ''
            , spread: true
            , children: [{
                title: '三级2-1-1'
                , id: 11
                , field: ''
            }, {
                title: '三级2-1-2'
                , id: 12
                , field: ''
            }]
        }, {
            title: '二级2-2'
            , id: 6
            , field: ''
            , children: [{
                title: '三级2-2-1'
                , id: 13
                , field: ''
            }, {
                title: '三级2-2-2'
                , id: 14
                , field: ''
                , disabled: true
            }]
        }]
    }, {
        title: '一级3'
        , id: 16
        , field: ''
        , children: [{
            title: '二级3-1'
            , id: 17
            , field: ''
            , fixed: true
            , children: [{
                title: '三级3-1-1'
                , id: 18
                , field: ''
            }, {
                title: '三级3-1-2'
                , id: 19
                , field: ''
            }]
        }, {
            title: '二级3-2'
            , id: 27
            , field: ''
            , children: [{
                title: '三级3-2-1'
                , id: 28
                , field: ''
            }, {
                title: '三级3-2-2'
                , id: 29
                , field: ''
            }]
        }]
    }]

    $.ajax({
        async: true,
        url: '/api/ssh/' + activeSshId + '/folder',
        type: 'GET',
        data: { page: 1 },
        dataType: 'json',
        timeout: 30000,
        success: successCallback,
        error: errorCallback,
        complete: completeCallback
    })

    function successCallback(res) {
        if (res.code == 0) {
            resData = res.data;
            let treeData = [];
            resData.list.forEach(item => {
                let treeItem = {
                    title: item.Name,
                    id: item.Id,
                    // field: item.name,
                    // spread: true,
                    children: [
                        {
                            title: '...',
                            id: item.Id,
                        }
                    ]
                }
                treeData.push(treeItem)
            });
            tree.render({
                elem: '#folder_tree'
                , data: treeData
                // , showCheckbox: true  //是否显示复选框
                // , id: 'demoId1'
                // , isJump: true //是否允许点击节点时弹出新窗口跳转
                , click: function (obj) {
                    var data = obj.data;  //获取当前点击的节点数据
                    layer.msg('状态：' + obj.state + '<br>节点数据：' + JSON.stringify(data));
                }
            });
        } else {
            layer.msg(res.msg, { icon: 2, time: 2000 });
        }
    }

    function errorCallback(err, status) {
        console.error("err:", err)
    }

    function completeCallback(xhr, status) {
    }



    // util.event('lay-demo', {
    //     getChecked: function (othis) {
    //         var checkedData = tree.getChecked('demoId1'); //获取选中节点的数据

    //         layer.alert(JSON.stringify(checkedData), { shade: 0 });
    //         console.log(checkedData);
    //     }
    //     , setChecked: function () {
    //         tree.setChecked('demoId1', [12, 16]); //勾选指定节点
    //     }
    //     , reload: function () {
    //         //重载实例
    //         tree.reload('demoId1', {

    //         });

    //     }
    // });
    //#endregion
});