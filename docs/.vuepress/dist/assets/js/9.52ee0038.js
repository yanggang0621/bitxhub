(window.webpackJsonp=window.webpackJsonp||[]).push([[9],{355:function(s,a,t){"use strict";t.r(a);var e=t(43),n=Object(e.a)({},(function(){var s=this,a=s.$createElement,t=s._self._c||a;return t("ContentSlotsDistributor",{attrs:{"slot-key":s.$parent.slotKey}},[t("p",[s._v("我们提供了脚本来快速启动两条Fabric应用链A和B，跨链网关和中继链。")]),s._v(" "),t("h2",{attrs:{id:"部署fabric网络"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#部署fabric网络"}},[s._v("#")]),s._v(" 部署Fabric网络")]),s._v(" "),t("p",[s._v("在运行跨链网络之前，必要的软件如Golang和Docker可以根据官网自行安装。确保 $GAPTH， $GOBIN等环境变量已经正确设置。")]),s._v(" "),t("p",[s._v("以上的软件依赖安装之后，我们提供了脚本来安装启动两个简单的Fabric网络。")]),s._v(" "),t("p",[s._v("下载部署fabric网络脚本（ffn.sh）：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[t("span",{pre:!0,attrs:{class:"token function"}},[s._v("wget")]),s._v(" https://github.com/meshplus/goduck/raw/master/scripts/quick_start/ffn.sh\n")])])]),t("p",[s._v("启动fabric网络：")]),s._v(" "),t("p",[t("strong",[s._v("注意：")]),s._v(" 脚本运行过程中按照提示进行确认即可。")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" ffn.sh down // 如果本地已经有Fabric网络运行，需要先关闭，如果没有可以不运行该命令\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" ffn.sh up //启动两条fabric应用链A和B\n")])])]),t("p",[t("strong",[s._v("注意：")]),s._v(" 命令执行完，会在当前目录生成"),t("code",[s._v("crypto-config")]),s._v("和"),t("code",[s._v("crypto-configB")]),s._v("文件夹，后面的"),t("code",[s._v("chaincode.sh")]),s._v("和 "),t("code",[s._v("fabric_pier.sh")]),s._v("需要在执行目录下存在上述文件。")]),s._v(" "),t("h2",{attrs:{id:"部署跨链合约"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#部署跨链合约"}},[s._v("#")]),s._v(" 部署跨链合约")]),s._v(" "),t("p",[s._v("下载操作"),t("code",[s._v("chaincode")]),s._v("的脚本（chaincode.sh）：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[t("span",{pre:!0,attrs:{class:"token function"}},[s._v("wget")]),s._v(" https://raw.githubusercontent.com/meshplus/goduck/master/scripts/quick_start/chaincode.sh\n")])])]),t("p",[s._v("拷贝"),t("code",[s._v("crypto-config")]),s._v("和"),t("code",[s._v("crypto-configB")]),s._v("到当前目录，执行以下命令：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[s._v("// -c：指定fabric cli的配置文件，默认为config.yaml\n// 应用链A部署chaincode\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" chaincode.sh "),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("install")]),s._v(" \n\n//应用链B部署chaincode\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" chaincode.sh "),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("install")]),s._v(" -c "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'configB.yaml'")]),s._v(" \n")])])]),t("p",[s._v("该命令会在指定fabric网络中部署"),t("code",[s._v("broker")]),s._v(", "),t("code",[s._v("transfer")]),s._v(" 和"),t("code",[s._v("data_swapper")]),s._v("三个"),t("code",[s._v("chaincode")])]),s._v(" "),t("p",[s._v("部署完成后，通过以下命令检查是否部署成功：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[s._v("// 上一步的脚本默认初始化了一个有10000余额的账号Alice（transfer chaincode）\n// -c：指定fabric cli的config.yaml配置文件\n// 查看应用链A中Alice的余额\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" chaincode.sh get_balance -c "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'config.yaml'")]),s._v(" \n****************************************************************************************************\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("0")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer0.org2.example.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("10000")]),s._v("\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("1")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer1.org2.example.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("10000")]),s._v("\n****************************************************************************************************\n\n//查看应用链B中Alice的余额\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" chaincode.sh get_balance -c "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'configB.yaml'")]),s._v(" \n****************************************************************************************************\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("0")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer0.org2.example1.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("10000")]),s._v("\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("1")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer1.org2.example1.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("10000")]),s._v("\n****************************************************************************************************\n")])])]),t("h1",{attrs:{id:"启动bitxhub"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#启动bitxhub"}},[s._v("#")]),s._v(" 启动BitXHub")]),s._v(" "),t("p",[s._v("BitXHub依赖于"),t("a",{attrs:{href:"https://golang.org/",target:"_blank",rel:"noopener noreferrer"}},[s._v("golang"),t("OutboundLink")],1),s._v(" 和"),t("a",{attrs:{href:"https://github.com/tmux/tmux/wiki",target:"_blank",rel:"noopener noreferrer"}},[s._v("tmux"),t("OutboundLink")],1),s._v("，需要提前进行安装。")]),s._v(" "),t("p",[s._v("使用下面的命令克隆项目：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[t("span",{pre:!0,attrs:{class:"token function"}},[s._v("git")]),s._v(" clone git@github.com:meshplus/bitxhub.git\n")])])]),t("p",[s._v("BitXHub还依赖于一些小工具，使用下面的命令进行安装：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[t("span",{pre:!0,attrs:{class:"token builtin class-name"}},[s._v("cd")]),s._v(" bitxhub\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("git")]),s._v(" checkout c3d850d3d4cb8742310bb8e2614e18cb075a3dc1\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" scripts/prepare.sh \n")])])]),t("p",[s._v("最后，运行下面的命令即可运行一个四节点的BitXHub中继链：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[t("span",{pre:!0,attrs:{class:"token function"}},[s._v("make")]),s._v(" cluster\n")])])]),t("p",[t("strong",[s._v("注意：")]),s._v(" "),t("code",[s._v("make cluster")]),s._v("启动会使用"),t("code",[s._v("tmux")]),s._v("进行分屏，所以在命令执行过程中，最好不要进行终端切换。")]),s._v(" "),t("h2",{attrs:{id:"启动跨链网关"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#启动跨链网关"}},[s._v("#")]),s._v(" 启动跨链网关")]),s._v(" "),t("p",[s._v("下载相关的脚本（fabric_pier.sh）：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[t("span",{pre:!0,attrs:{class:"token function"}},[s._v("wget")]),s._v(" https://github.com/meshplus/goduck/raw/master/scripts/quick_start/fabric_pier.sh\n")])])]),t("p",[s._v("执行以下命令即可启动跨链网关：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[s._v("//启动跨链网关连接应用链A和BitXHub\n// -r: 跨链网关启动目录，默认为.pier目录\n// -c: fabric组织证书目录，默认为crypto-config\n// -g: 指定fabric cli连接的配置文件，默认为config.yaml\n// -p: 跨链网关的启动端口，默认为8987\n// -b: 中继链GRPC地址，默认为localhost:60011\n// -o: pprof端口，默认为44555\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" fabric_pier.sh start -r "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'.pier'")]),s._v(" -c "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'crypto-config'")]),s._v("  -g "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'config.yaml'")]),s._v(" -p "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("8987")]),s._v(" -b "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'localhost:60011'")]),s._v(" -o "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("44555")]),s._v("\n\n//启动跨链网关连接应用链B和BitXHub\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" fabric_pier.sh start -r "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'.pierB'")]),s._v(" -c "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'crypto-configB'")]),s._v("  -g "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'configB.yaml'")]),s._v(" -p "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("8988")]),s._v(" -b "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'localhost:60011'")]),s._v(" -o "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("44556")]),s._v("\n")])])]),t("p",[s._v("在该目录下通过以下命令可以得到该跨链网关对应应用链的ID：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[s._v("//应用链A的ID\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" fabric_pier.sh "),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("id")]),s._v(" -r "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'.pier'")]),s._v("\n\n//应用链B的ID\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" fabric_pier.sh "),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("id")]),s._v(" -r "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'.pierB'")]),s._v("\n")])])]),t("p",[t("strong",[s._v("注意：")]),s._v(" 后面跨链交易命令需要该值")]),s._v(" "),t("h2",{attrs:{id:"跨链转账"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#跨链转账"}},[s._v("#")]),s._v(" 跨链转账")]),s._v(" "),t("p",[s._v("使用"),t("strong",[s._v("部署跨链合约")]),s._v("章节中下载的"),t("code",[s._v("chaincode.sh")]),s._v("进行相关"),t("code",[s._v("chaincode")]),s._v("调用。")]),s._v(" "),t("ol",[t("li",[s._v("查询Alice余额")])]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[s._v("// 查询应用链A中Alice的余额\n// -c：指定fabric cli的config.yaml配置文件\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" chaincode.sh get_balance -c "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'config.yaml'")]),s._v("\n****************************************************************************************************\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("0")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer0.org2.example.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("10000")]),s._v("\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("1")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer1.org2.example.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("10000")]),s._v("\n****************************************************************************************************\n\n// 查询应用链B中Alice的余额\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" chaincode.sh get_balance -c "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'configB.yaml'")]),s._v("\n****************************************************************************************************\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("0")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer0.org2.example1.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("10000")]),s._v("\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("1")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer1.org2.example1.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("10000")]),s._v("\n****************************************************************************************************\n")])])]),t("ol",{attrs:{start:"2"}},[t("li",[s._v("发送一笔跨链转账")])]),s._v(" "),t("p",[s._v("下面的命令会将应用链A中Alice的一块钱转移到应用链B中Alice：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[s._v("// -c：指定fabric cli的config.yaml配置文件\n// -t: 目的链的ID（应用链B的ID）\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" chaincode.sh interchain_transfer -c "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'config.yaml'")]),s._v(" -t "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("<")]),s._v("target_appchain_id"),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v(">")]),s._v("\n")])])]),t("ol",{attrs:{start:"3"}},[t("li",[s._v("查询余额")])]),s._v(" "),t("p",[s._v("分别在两条链上查询Alice余额：")]),s._v(" "),t("div",{staticClass:"language-shell extra-class"},[t("pre",{pre:!0,attrs:{class:"language-shell"}},[t("code",[s._v("// 查询应用链A中Alice的余额，发现余额少了一块钱\n// -c：指定fabric cli的config.yaml配置文件\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" chaincode.sh get_balance -c "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'config.yaml'")]),s._v("\n****************************************************************************************************\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("0")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer0.org2.example.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("9999")]),s._v("\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("1")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer1.org2.example.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("9999")]),s._v("\n****************************************************************************************************\n\n// 查询应用链B中Alice的余额，发现余额多了一块钱\n"),t("span",{pre:!0,attrs:{class:"token function"}},[s._v("bash")]),s._v(" chaincode.sh get_balance -c "),t("span",{pre:!0,attrs:{class:"token string"}},[s._v("'configB.yaml'")]),s._v("\n****************************************************************************************************\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("0")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer0.org2.example1.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("10001")]),s._v("\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Response"),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("[")]),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("1")]),t("span",{pre:!0,attrs:{class:"token punctuation"}},[s._v("]")]),s._v(": //peer1.org2.example1.com\n***** "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  "),t("span",{pre:!0,attrs:{class:"token operator"}},[s._v("|")]),s._v("  Payload: "),t("span",{pre:!0,attrs:{class:"token number"}},[s._v("10001")]),s._v("\n****************************************************************************************************\n")])])]),t("p",[t("strong",[s._v("注意：")]),s._v(" "),t("code",[s._v("chaincode.sh")]),s._v("调用不同的fabric，需要不同的"),t("code",[s._v("config.yaml")]),s._v("，注意区分。")])])}),[],!1,null,null,null);a.default=n.exports}}]);