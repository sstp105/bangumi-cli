package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/handler/formathandler"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "格式化所有媒体文件，使其符合 Plex / Jellyfin / Emby 媒体流格式",
	Long: `
Summary:
  rename 命令将通过将标题、季节和集数附加到每个文件名来重命名当前目录中的视频文件。
  它将处理常见的媒体格式文件（如 .mp4、.mkv），并递归地处理子目录。
  季节和集数的编号将根据文件夹的名称来确定，文件名将相应地重命名。
  如果文件夹包含自定义格式的季节，例如 夏目友人帐 柒，代表第七季，建议将文件夹重命名为
  夏目友人帐 柒 - S07、夏目友人帐 柒 第七季，或者 夏目友人帐 第七期。
  其他示例：东京喰种√A，CLANNAD 〜AFTER STORY〜，Code Geass 反叛的鲁路修R2.
	`,
	Example: `
  bangumi format 默认将处理当前工作目录下的所有媒体文件。
`,
	Run: func(cmd *cobra.Command, args []string) {
		formathandler.Run()
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
}
