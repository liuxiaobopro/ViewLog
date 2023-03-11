layui.use(['tree', 'code', 'dropdown'], function () {
    layui.code()
    var tree = layui.tree,
        layer = layui.layer,
        $ = layui.jquery,
        form = layui.form,
        util = layui.util

    //#region 获取隐藏域
    var activeSshId = $("#activeSshId").val(); // 当前激活的sshId
    console.log("activeSshId:", activeSshId);

    var activeFolderId = 0; // 当前激活的文件夹Id
    var selectFilename = ""; // 当前选中的文件名
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
                layer.msg(res.msg, { icon: 2, time: 2000 }, function () {
                    location.reload();
                });
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
        if (activeSshId == 0) {
            layer.msg("请先选择某个SSH", { icon: 2, time: 2000 });
            return false;
        }
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

    //#region 监听文件夹切换
    form.on('select(folder-list)', function (data) {
        activeFolderId = data.value;
        console.log("activeFolderId:", activeFolderId);
        var html = ""
        if (activeFolderId > 0) {
            var loading = layer.load(3)
            $.ajax({
                async: true,
                url: '/api/folder/' + activeFolderId + '/child',
                type: 'GET',
                data: {},
                dataType: 'json',
                timeout: 30000,
                success: successCallback,
                error: errorCallback,
                complete: completeCallback
            })

            function successCallback(res) {
                console.log("res:", res);
                if (res.code == 0) {
                    var data = res.data;
                    if (data.length > 0) {
                        data.forEach(item => {
                            html += `
                            <div class="file-item">${item}</div>
                            `
                        });
                    } else {
                        html = "该文件夹下没有文件"
                    }
                    $(".file-list").html(html);
                } else {
                    layer.msg(res.msg, { icon: 2, time: 2000 });
                }
            }

            function errorCallback(err, status) {
                console.error("err:", err)
            }

            function completeCallback(xhr, status) {
                console.log('Ajax请求已结束。');
                layer.close(loading)
            }
        } else {
            html = "请选择文件夹"
        }
        $(".file-list").html(html);
    });
    //#endregion

    //#region 点击文件列表
    $(".file-list").on("click", ".file-item", function () {
        var fileName = $(this).text();
        selectFilename = fileName;
        $("#file-name").html(selectFilename);
        getFileDetail()
    })
    //#endregion

    //#region FUNC 获取文件详情
    function getFileDetail() {
        var loading = layer.load(3)
        jsonData = {
            "folderId": activeFolderId,
            "path": selectFilename
        }
        $.ajax({
            async: true,
            url: '/api/file',
            type: 'GET',
            data: jsonData,
            dataType: 'json',
            timeout: 30000,
            success: successCallback,
            error: errorCallback,
            complete: completeCallback
        })

        function successCallback(res) {
            console.log("res:", res);
            if (res.code == 0) {
                $("#content").html(res.data);
                setTimeout(() => {
                    var h4 = $('.out-box').prop("scrollHeight"); //等同 $('.out-box')[0].scrollHeight
                    $('.out-box').scrollTop(h4);
                }, 10);
            } else {
                layer.msg(res.msg, { icon: 2, time: 2000 });
            }
        }

        function errorCallback(err, status) {
            console.error("err:", err)
        }

        function completeCallback(xhr, status) {
            console.log('Ajax请求已结束。');
            layer.close(loading)
        }
    }
    //#endregion

    //#region 重新获取
    $("#getDetailAgain").click(function () {
        getFileDetail()
    })
    //#endregion
});