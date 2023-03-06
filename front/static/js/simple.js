layui.use(['tree', 'code'], function () {
    layui.code()
    var tree = layui.tree, layer = layui.layer, $ = layui.jquery, form = layui.form

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
            area: ['50%','50%'],
            content: 'folder_add'
        });
    })
    //#endregion

});