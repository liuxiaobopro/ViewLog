layui.use(['layer', 'form'], function () {
    var layer = layui.layer,
        form = layui.form,
        $ = layui.jquery;

    form.on('submit(install)', function (data) {
        let loading = layer.load(3)
        console.log("data:", data);
        let f = data.field;
        let reqData = {
            host: f.host,
            port: parseInt(f.port),
            dbname: f.dbname,
            user: f.user,
            password: f.password,
            charset: f.charset,
        }
        $.ajax({
            url: '/install',
            type: 'POST',
            data: JSON.stringify(reqData),
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
                    window.location.href = "/";
                });
            } else {
                layer.msg(res.msg, { icon: 2, time: 1000 });
            }
        }

        function errorCallback(err, status) {
            console.error("err:", err)
        }

        function completeCallback(xhr, status) {
            console.log('Ajax请求已结束。');
            layer.close(loading);
        }
        return false;
    });
});