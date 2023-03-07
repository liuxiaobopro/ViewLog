layui.use(['layer', 'form'], function () {
    var layer = layui.layer,
        form = layui.form,
        $ = layui.jquery;

    $("#reset").click(function () {
        layer.confirm('确定要重置吗？', { icon: 3, title: '提示' }, function (index) {
            layer.close(index);
            let loading = layer.load(3);
            $.ajax({
                async: true,
                url: '/reset',
                type: 'POST',
                dataType: 'json',
                timeout: 30000,
                success: successCallback,
                error: errorCallback,
                complete: completeCallback
            })

            function successCallback(res) {
                console.log("res:", res);
                if (res.code == 0) {
                    layer.msg(res.msg, { icon: 6, time: 1000 }, function () {
                        location.reload();
                    });
                } else {
                    layer.msg(res.msg, { icon: 5, time: 1000 });
                }
            }

            function errorCallback(err, status) {
                console.error("err:", err)
            }

            function completeCallback(xhr, status) {
                console.log('Ajax请求已结束。');
                layer.close(loading);
            }
        });
    })
})