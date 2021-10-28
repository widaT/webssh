<script>
import 'xterm/css/xterm.css'
import "./xterm.css";
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import Base64 from "crypto-js/enc-base64";
import Utf8 from "crypto-js/enc-utf8";
const msgData = '1';
const msgResize = '2';
export default {
    mounted() {
        const terminal = new Terminal({})
        const fitAddon = new FitAddon()
        terminal.loadAddon(fitAddon)
        fitAddon.fit()
        let terminalContainer = document.getElementById("app");
        const webSocket = new WebSocket(`ws://127.0.0.1:8080/ws/1`);

        webSocket.onmessage = (event) => {
            terminal.write(event.data.toString(Utf8));
        };

        webSocket.onopen = () => {
            terminal.open(terminalContainer);
            fitAddon.fit();
            terminal.write("welcome to WebSSH ☺\r\n");
            terminal.focus();
        };

        webSocket.onclose = () => {
            terminal.write("\r\nWebSSH quit!");
        };

        webSocket.onerror = (event) => {
            console.error(event);
            webSocket.close();
        };

        terminal.onKey((event) => {
            webSocket.send(msgData + Base64.stringify(Utf8.parse(event.key)));
        });

        terminal.onResize(({ cols, rows }) => {
            console.log(cols,rows)
            webSocket.send(msgResize +
                Base64.stringify(
                    Utf8.parse(
                        JSON.stringify({
                            columns: cols,
                            rows: rows
                        })
                    )
                )
            );
        });
        // 内容全屏显示-窗口大小发生改变时
        // resizeScreen
        window.addEventListener("resize", () =>{
            fitAddon.fit();
        },false);

    },
    methods: {
    },
}
</script>