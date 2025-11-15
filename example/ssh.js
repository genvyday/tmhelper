// 自动登录ssh，并停留在交互式shell

tp=tmh.Exec(["ssh", "xr@127.0.0.1"]); // 运行命令

// 执行多个匹配，默认命中任意一个就返回
tp.Matchs([
    ["yes/no", "yes\n", "C"],       // "C"(continue)标志表示命中后不退出，继续匹配
    ["password", "ya8j3fpa*ed\n"],
]);
tp.Matchs([["$", "cd /data/git/\n"]]);

tp.Term(); // 停留在交互式终端，若直接结束则可调用 tp.Exit()