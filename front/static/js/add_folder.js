layui.use(['layer', 'form'], function () {
    var form = layui.form, $ = layui.jquery;

    form.on('submit(submit)', function (data) {
        let formData = data.field;
        console.log("formData:", formData);
        if (formData.sshId == "") {
            layer.msg("请选择SSH", { icon: 2, time: 2000 });
            return false;
        }
        if (formData.name == "") {
            layer.msg("名称不能为空", { icon: 2, time: 2000 });
            return false;
        }
        if (formData.path == "") {
            layer.msg("路径不能为空", { icon: 2, time: 2000 });
            return false;
        }

        let reqData = {
            sshId: parseInt(formData.sshId),
            name: formData.name,
            path: formData.path,
        }
        $.ajax({
            url: '/api/folder',
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
            console.log('Ajax请求已结束。');
        }
        return false;
    });
});