# Changelog

## Version 1.0.0 (2025-11-28)

### Added

#### Licensing
- Added AGPL-3.0 license file
- Added AGPL-3.0 license headers to all Go source files
- Updated README with comprehensive license information
- Added explanation of AGPL-3.0 network use requirements

#### Commands
- Added `!source` command (aliases: `!src`, `!info`)
  - Displays GitHub repository link
  - Shows creator information (GitHub: blubskye, Discord: blubaustin)
  - Shows license information (AGPL-3.0)
  - Shows version number
  - Satisfies AGPL-3.0 source code availability requirement

#### Command Line Flags
- Added `--version` flag
  - Shows bot version (1.0.0)
  - Shows author information
  - Shows source repository URL
  - Shows license type
- Added `--trace` flag
  - Enables full stack tracing for debugging
  - Uses `runtime/debug.SetTraceback("all")`
  - Helpful for diagnosing crashes and errors

#### Visual Features
- Added ASCII art display on bot startup
  - Displays beautiful Braille/Unicode art from `ascii.txt`
  - Shows when bot successfully connects to Discord
  - Includes guild count in startup message

#### Documentation
- Added "Command Line Flags" section to README
- Updated command tables to include `!source` command
- Added creator information section
- Enhanced license section with AGPL-3.0 compliance details
- Added copyright notices

### Changed
- Updated all copyright years from 2024 to 2025
- Updated README license section from generic to AGPL-3.0 specific
- Enhanced help command to include `!source` command
- Updated support section with GitHub issue tracker link

### Technical Details

**Files Modified:**
- `LICENSE` - Added AGPL-3.0 full license text
- `ascii.txt` - ASCII art for startup display
- `cmd/bot/main.go` - Added flags, version info, trace support
- `internal/bot/bot.go` - Added license header, ASCII art display function
- `internal/bot/config.go` - Added license header
- `internal/commands/commands.go` - Added license header, source command
- `internal/database/database.go` - Added license header
- `internal/music/player.go` - Added license header
- `internal/permissions/permissions.go` - Added license header
- `internal/queue/queue.go` - Added license header
- `README.md` - Added flags documentation, license info, creator info
- `CHANGELOG.md` - Documentation of all changes

**New Constants:**
- `Version = "1.0.0"`
- `Author = "blubskye (blubaustin)"`
- `Repo = "https://github.com/blubskye/miku_discord_music_bot"`

**AGPL-3.0 Compliance:**
All source files now include proper license headers as required by AGPL-3.0, and the `!source` command provides users with easy access to the source code repository, satisfying the network use provision of the license.
