layui.use(['tree', 'code', 'dropdown'], function () {
    layui.code()
    var tree = layui.tree,
        layer = layui.layer,
        $ = layui.jquery,
        form = layui.form

    //#region 监听运行模式
    form.on('radio(mode)', function (data) {
        console.log("data:", data);
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
            url: '/api/ssh',
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
            }else{
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
});