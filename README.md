# ClearCode

This is a file editor I am making that would work like nvim. Currently, the only features are Saving / Opening files, and writing to them.

I included 3 different modes: `FILE`, `NORMAL`, `INSERT`.

File Mode allows you to browse the current file directory using h, j, k, l. You go into file mode by pressing your f key, and exiting by pressing your Escape Key.

Normal Mode (when a file is selected) moves a cursor around in a file. Using h, j, k, l like nvim.

Currently the only keybinds for normal mode is:
- ctrl+c, q: QUIT
- $: Jump to end of the line
- shift+A: Jump to end of the line in insert mode
- 0: Jump to start of the line
- ctrl+s: Save File

From normal mode, you can enter `FILE` and `INSERT` mode by pressing f or i.

## Running

This is only runnable on linux.

Install the binary from the current Release Page.
You need to run this command to make sure that the executable can run: `chmod +x clearcode-darwin-arm64`.
After that, you can just run `./clearcode-linux-amd64`.
