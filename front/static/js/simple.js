layui.use(['tree', 'code'], function () {
    layui.code()
    var tree = layui.tree, layer = layui.layer, $ = layui.jquery, form = layui.form

    //#region 监听运行模式
    form.on('radio(mode)', function (data) {
        console.log("data:", data);
        var mode = data.value;
        if (mode == "1") {
            $(".do_search").hide();
        }else{
            $(".do_search").show();
        }
    });
    //#endregion
});