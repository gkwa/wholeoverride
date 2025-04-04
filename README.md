# WholeOverride

WholeOverride is a command-line tool that generates a consolidated markdown index from a collection of recipe markdown files. It helps organize and index recipe collections by creating a navigable table of contents and providing a clear view of recipes and their creators.

## What It Does

1. It scans a directory for markdown files with recipe information
2. It parses recipe files that have a "filetype: recipe" in their frontmatter
3. It extracts metadata including recipe title, image URL, and creator information
4. It also processes separate creator files to get their information
5. It generates an index markdown file in either "sections" or "table" format
6. It creates a table of contents for easy navigation

## Installation

```bash
# Clone the repository
git clone https://github.com/gkwa/wholeoverride.git

# Navigate to the project directory
cd wholeoverride

# Build the application
go build -o wholeoverride
```

## Usage

### Generate Command

The primary functionality is provided by the `generate` command:

```bash
./wholeoverride generate --basedir /path/to/recipes --format sections
```

Options:

- `--basedir`: (Required) Path to the directory containing recipe markdown files
- `--format`: (Optional) Output format - either "sections" or "table" (default: "sections")

### Version Command

Display version information:

```bash
./wholeoverride version
```

### Global Options

These options can be used with any command:

- `--verbose, -v`: Increase output verbosity (can be used multiple times for more detail)
- `--log-format`: Set log format to "json" or "text" (default: "text")
- `--config`: Specify a config file path (default: $HOME/.wholeoverride.yaml)

## File Structure Requirements

### Recipe Files

Recipe files should be markdown files with frontmatter that includes:

```yaml
---
filetype: recipe
pic: "image-file-path.jpg" # or a URL
creator: "Creator Name"
---
Recipe content goes here...
```

### Creator Files

Creator files should be markdown files named after the creator (e.g., `Creator Name.md`) with frontmatter:

```yaml
---
pic: "creator-image-path.jpg" # or a URL
---
Creator information goes here...
```

## Output

The tool generates a file called `recipeindex.md` in the base directory with:

1. A table of contents for all recipes
2. Recipe sections with links to both recipe and creator files
3. Images for both recipes and creators
4. Navigation links back to the table of contents

## Example Output (Sections Format)

```markdown
# TOC

- [[#Apple Pie|Apple Pie]] ^apple-pie
- [[#Chocolate Cake|Chocolate Cake]] ^chocolate-cake

## Apple Pie

[[#^apple-pie|toc]]

| [[Apple Pie]]               | [[Jane Baker]]      |
| --------------------------- | ------------------- |
| ![Apple Pie](apple-pie.jpg) | ![[jane-baker.jpg]] |

## Chocolate Cake

[[#^chocolate-cake|toc]]

| [[Chocolate Cake]]               | [[John Chef]]      |
| -------------------------------- | ------------------ |
| ![Chocolate Cake](choc-cake.jpg) | ![[john-chef.jpg]] |
```

## How It Works

1. **Scanning**: The tool walks through the specified directory to find all markdown files.
2. **Parsing**: It reads each file and parses frontmatter using the Goldmark library.
3. **Filtering**: Files with `filetype: recipe` are identified as recipes.
4. **Creator Lookup**: For each recipe, the tool finds the corresponding creator file.
5. **Slug Generation**: For each recipe, a slug is generated for linking purposes.
6. **Content Generation**: Based on the chosen format, the tool generates markdown content.
7. **TOC Creation**: A table of contents is generated with links to each recipe.
8. **File Writing**: The final content is written to `recipeindex.md`.

## Logging

The tool provides extensive logging capabilities:

- Level 0 (default): Info level, basic information
- Level 1 (-v): Debug level, processing details
- Level 2 (-vv): Trace level, detailed file processing information
