layui.use(['form'], function () {
    var form = layui.form, $ = layui.jquery, layer = layui.layer;
    form.on('submit(submit)', function (data) {
        let loading = layer.load(3, { shade: [0.3, '#fff'] });
        formData = data.field;
        if (formData.name == "") {
            layer.msg("主机名称不能为空", { icon: 2, time: 2000 });
            layer.close(loading);
            return false;
        }
        if (formData.host == "") {
            layer.msg("主机地址不能为空", { icon: 2, time: 2000 });
            layer.close(loading);
            return false;
        }
        if (formData.port == "") {
            layer.msg("端口不能为空", { icon: 2, time: 2000 });
            layer.close(loading);
            return false;
        }

        jsonStr = JSON.stringify({
            "name": formData.name,
            "host": formData.host,
            "port": parseInt(formData.port),
            "username": formData.username,
            "password": formData.password,
        });

        $.ajax({
            async: true,
            url: '/api/ssh',
            type: 'POST',
            data: jsonStr,
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
                    // 关闭当前窗口
                    var index = parent.layer.getFrameIndex(window.name);
                    parent.layer.close(index);
                    // 刷新父窗口
                    parent.location.reload();
                });
            } else {
                layer.msg(res.msg, { icon: 2, time: 2000 });
            }
        }

        function errorCallback(err, status) {
            console.error("err:", err)
        }

        function completeCallback(xhr, status) {
            layer.close(loading);
            console.log('Ajax请求已结束。');
        }
    })
})