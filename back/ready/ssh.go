package ready

import toolsSsh "ViewLog/back/tools/ssh"

func Ssh() {
	_ = toolsSsh.UpdateGlobalClient()
}
