#!/usr/bin/env bash
# verify-plugin-package.sh — Validate the Game Designer plugin package structure
#
# Usage: ./scripts/verify-plugin-package.sh
#
# Checks manifests, root skills, bundled assets, and documentation references.
# Run from the repository root.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

pass=0
fail=0
errors=()

check() {
  local name="$1"
  shift
  if "$@" 2>/dev/null; then
    echo -e "  ${GREEN}PASS${NC} $name"
    pass=$((pass + 1))
  else
    echo -e "  ${RED}FAIL${NC} $name"
    fail=$((fail + 1))
    errors+=("$name")
  fi
}

echo "=== Game Designer Plugin — Package Validation ==="
echo ""

# 1. Manifest JSON validity
echo "1. Manifest files"
check ".claude-plugin/plugin.json is valid JSON" python3 -c "import json; json.load(open('$ROOT_DIR/.claude-plugin/plugin.json'))"
check ".claude-plugin/marketplace.json is valid JSON" python3 -c "import json; json.load(open('$ROOT_DIR/.claude-plugin/marketplace.json'))"
check ".codex-plugin/plugin.json is valid JSON" python3 -c "import json; json.load(open('$ROOT_DIR/.codex-plugin/plugin.json'))"

# 2. Required manifest fields
echo ""
echo "2. Manifest fields"
check "Claude manifest has name" python3 -c "import json; d=json.load(open('$ROOT_DIR/.claude-plugin/plugin.json')); assert 'name' in d and d['name']"
check "Claude manifest has description" python3 -c "import json; d=json.load(open('$ROOT_DIR/.claude-plugin/plugin.json')); assert 'description' in d and d['description']"
check "Claude manifest has skills path" python3 -c "import json; d=json.load(open('$ROOT_DIR/.claude-plugin/plugin.json')); assert 'skills' in d and d['skills']"
check "Marketplace has owner" python3 -c "import json; d=json.load(open('$ROOT_DIR/.claude-plugin/marketplace.json')); assert 'owner' in d and d['owner']['name']"
check "Marketplace has plugins" python3 -c "import json; d=json.load(open('$ROOT_DIR/.claude-plugin/marketplace.json')); assert 'plugins' in d and len(d['plugins']) > 0"
check "Codex manifest has name" python3 -c "import json; d=json.load(open('$ROOT_DIR/.codex-plugin/plugin.json')); assert 'name' in d and d['name']"
check "Codex manifest has skills path" python3 -c "import json; d=json.load(open('$ROOT_DIR/.codex-plugin/plugin.json')); assert 'skills' in d and d['skills']"
check "Both manifests share same plugin name" python3 -c "
import json
c=json.load(open('$ROOT_DIR/.claude-plugin/plugin.json'))
d=json.load(open('$ROOT_DIR/.codex-plugin/plugin.json'))
assert c['name'] == d['name'], f'{c[\"name\"]} != {d[\"name\"]}'
"
check "Both manifests share same skills path" python3 -c "
import json
c=json.load(open('$ROOT_DIR/.claude-plugin/plugin.json'))
d=json.load(open('$ROOT_DIR/.codex-plugin/plugin.json'))
assert c['skills'] == d['skills'], f'{c[\"skills\"]} != {d[\"skills\"]}'
"

# 3. Skills discovery
echo ""
echo "3. Skills"
SKILLS_DIR="$ROOT_DIR/skills"
check "Skills directory exists" test -d "$SKILLS_DIR"
check "Deprecated plugin/skills directory absent" test ! -d "$ROOT_DIR/plugin/skills"
check "Deprecated plugins/game-designer wrapper absent" test ! -e "$ROOT_DIR/plugins/game-designer"
check "Optional .agents local catalog is valid when present" python3 -c "
import json, pathlib
path = pathlib.Path('$ROOT_DIR/.agents/plugins/marketplace.json')
if path.exists():
    data = json.load(path.open())
    plugin = next((p for p in data.get('plugins', []) if p.get('name') == 'game-designer'), None)
    assert plugin is not None
    assert plugin.get('source', {}).get('source') == 'local'
    assert plugin.get('source', {}).get('path') == '.'
"
check "Claude manifest points at ./skills/" python3 -c "import json; d=json.load(open('$ROOT_DIR/.claude-plugin/plugin.json')); assert d.get('skills') == './skills/'"
check "Codex manifest points at ./skills/" python3 -c "import json; d=json.load(open('$ROOT_DIR/.codex-plugin/plugin.json')); assert d.get('skills') == './skills/'"

# Enumerate skills
skill_count=0
skill_names=""
if [ -d "$SKILLS_DIR" ]; then
  for skill_dir in "$SKILLS_DIR"/*/; do
    skill_name=$(basename "$skill_dir")
    skill_file="$skill_dir/SKILL.md"
    if [ -f "$skill_file" ]; then
      echo -e "  ${GREEN}PASS${NC} skill '$skill_name' has SKILL.md"
      pass=$((pass + 1))

      # Check frontmatter
      if head -5 "$skill_file" | grep -q "^name:"; then
        echo -e "  ${GREEN}PASS${NC} skill '$skill_name' has name in frontmatter"
        pass=$((pass + 1))
      else
        echo -e "  ${RED}FAIL${NC} skill '$skill_name' missing name in frontmatter"
        fail=$((fail + 1))
        errors+=("skill-$skill_name-frontmatter-name")
      fi

      if head -5 "$skill_file" | grep -q "^description:"; then
        echo -e "  ${GREEN}PASS${NC} skill '$skill_name' has description in frontmatter"
        pass=$((pass + 1))
      else
        echo -e "  ${RED}FAIL${NC} skill '$skill_name' missing description in frontmatter"
        fail=$((fail + 1))
        errors+=("skill-$skill_name-frontmatter-description")
      fi

      skill_count=$((skill_count + 1))
      skill_names="$skill_names $skill_name"
    else
      echo -e "  ${RED}FAIL${NC} skill directory '$skill_name' missing SKILL.md"
      fail=$((fail + 1))
      errors+=("skill-$skill_name-missing-skill-md")
    fi
  done
fi

check "At least 6 skills found (got $skill_count)" test "$skill_count" -ge 6

# 4. Duplicate skill names
echo ""
echo "4. Skill uniqueness"
dup_check=$(echo "$skill_names" | tr ' ' '\n' | sort | uniq -d | tr -d ' ')
if [ -z "$dup_check" ]; then
  echo -e "  ${GREEN}PASS${NC} no duplicate skill names"
  pass=$((pass + 1))
else
  echo -e "  ${RED}FAIL${NC} duplicate skill names: $dup_check"
  fail=$((fail + 1))
  errors+=("duplicate-skills")
fi

# 5. Bundled assets
echo ""
echo "5. Bundled assets"
for dir in server-template cli sdk-js contracts examples scripts; do
  check "$dir/ exists at plugin root" test -d "$ROOT_DIR/$dir"
done

# 6. Documentation references
echo ""
echo "6. Documentation"
check "README.md exists" test -f "$ROOT_DIR/README.md"
check "docs/integration/plugin-installation.md exists" test -f "$ROOT_DIR/docs/integration/plugin-installation.md"
check "docs/integration/agent-golden-path.md exists" test -f "$ROOT_DIR/docs/integration/agent-golden-path.md"
check "cli/README.md references setup skill" grep -q "gd-setup-cli" "$ROOT_DIR/cli/README.md"
check "README references root skills directory" grep -q "skills/" "$ROOT_DIR/README.md"
check "Install docs reference root skills directory" grep -q "skills/" "$ROOT_DIR/docs/integration/plugin-installation.md"
check "User-facing docs do not advertise deprecated plugin roots" python3 -c "
from pathlib import Path
files = [
    Path('$ROOT_DIR/README.md'),
    Path('$ROOT_DIR/docs/integration/plugin-installation.md'),
    Path('$ROOT_DIR/docs/integration/agent-golden-path.md'),
    Path('$ROOT_DIR/docs/integration/local-verification.md'),
]
bad = []
for path in files:
    text = path.read_text()
    for marker in ['plugin/skills/', 'plugins/game-designer/']:
        if marker in text:
            bad.append(f'{path}:{marker}')
assert not bad, ', '.join(bad)
"

# 7. Install docs do not claim CLI is built during install
echo ""
echo "7. Install docs accuracy"
check "Install docs do not claim auto-CLI-build" python3 -c "
text = open('$ROOT_DIR/docs/integration/plugin-installation.md').read().lower().replace('*', '')
assert 'does not compile' in text.lower() or 'does not build' in text.lower()
"

# 8. Slot machine theme consistency
echo ""
echo "8. Slot machine theme"
check "Skills do not reference stale activity-game example" python3 -c "
from pathlib import Path
skills_dir = Path('$ROOT_DIR/skills')
bad = []
for skill_file in skills_dir.glob('*/SKILL.md'):
    text = skill_file.read_text()
    if 'basic-activity-game' in text:
        bad.append(f'{skill_file}:basic-activity-game')
assert not bad, ', '.join(bad)
"
check "Skills reference slot machine concepts" python3 -c "
from pathlib import Path
skills_dir = Path('$ROOT_DIR/skills')
connect_skill = skills_dir / 'gd-connect-sdk' / 'SKILL.md'
text = connect_skill.read_text()
assert 'spin' in text.lower() and 'balance' in text.lower(), 'gd-connect-sdk skill must reference spin and balance'
"
check "Example h5-slot-machine exists" test -d "$ROOT_DIR/examples/h5-slot-machine"
check "Old example h5-activity-game absent" test ! -d "$ROOT_DIR/examples/h5-activity-game"

# 9. Skill naming convention (gd- prefix)
echo ""
echo "9. Skill naming convention"
check "All skill names start with gd- prefix" python3 -c "
from pathlib import Path
skills_dir = Path('$ROOT_DIR/skills')
bad = []
for skill_dir in sorted(skills_dir.iterdir()):
    if skill_dir.is_dir():
        if not skill_dir.name.startswith('gd-'):
            bad.append(skill_dir.name)
assert not bad, f'Skills missing gd- prefix: {bad}'
"

# 10. Skill directory matches frontmatter name
echo ""
echo "10. Skill directory/frontmatter alignment"
check "Every skill directory name matches its frontmatter name" python3 -c "
import re
from pathlib import Path
skills_dir = Path('$ROOT_DIR/skills')
bad = []
for skill_dir in sorted(skills_dir.iterdir()):
    if skill_dir.is_dir():
        skill_file = skill_dir / 'SKILL.md'
        if skill_file.exists():
            text = skill_file.read_text()
            m = re.search(r'^name:\s*(.+)$', text, re.MULTILINE)
            if m:
                frontmatter_name = m.group(1).strip()
                if frontmatter_name != skill_dir.name:
                    bad.append(f'{skill_dir.name} vs {frontmatter_name}')
            else:
                bad.append(f'{skill_dir.name}: no name in frontmatter')
assert not bad, f'Mismatches: {bad}'
"

# 11. Stale old skill name scan (active surfaces only)
echo ""
echo "11. Stale skill name scan"
check "Active surfaces contain no old skill names" python3 -c "
from pathlib import Path
old_names = ['setup-game-designer-cli', 'create-game-server', 'connect-js-sdk', 'deploy-game-server', 'debug-server-integration']
active_files = [
    Path('$ROOT_DIR/.codex-plugin/plugin.json'),
    Path('$ROOT_DIR/README.md'),
    Path('$ROOT_DIR/cli/README.md'),
    Path('$ROOT_DIR/docs/integration/plugin-installation.md'),
    Path('$ROOT_DIR/docs/integration/agent-golden-path.md'),
    Path('$ROOT_DIR/docs/integration/local-verification.md'),
]
# Also scan all skill files
skills_dir = Path('$ROOT_DIR/skills')
for skill_file in sorted(skills_dir.glob('*/SKILL.md')):
    active_files.append(skill_file)
bad = []
for path in active_files:
    if not path.exists():
        continue
    text = path.read_text()
    for old_name in old_names:
        if old_name in text:
            bad.append(f'{path.relative_to(Path(\"$ROOT_DIR\"))}:{old_name}')
assert not bad, f'Stale references: {bad}'
"

# 12. Deploy docs must not claim fake is the only provider
echo ""
echo "12. Production provider documentation"
check "Deploy docs reference production provider" python3 -c "
from pathlib import Path
text = Path('$ROOT_DIR/docs/deployment/paas-provider.md').read_text()
assert '3os' in text, 'paas-provider.md must reference the 3os production provider'
assert 'fake' in text, 'paas-provider.md should still document the fake provider'
"
check "CLI README documents production provider" python3 -c "
from pathlib import Path
text = Path('$ROOT_DIR/cli/README.md').read_text()
assert '3os' in text, 'cli/README.md must document the 3os provider'
assert 'fake' in text, 'cli/README.md must still document the fake provider'
"
check "Deploy skill documents both providers" python3 -c "
from pathlib import Path
text = Path('$ROOT_DIR/skills/gd-deploy-server/SKILL.md').read_text()
assert '3os' in text, 'gd-deploy-server skill must reference 3os provider'
assert 'fake' in text, 'gd-deploy-server skill must still reference fake provider'
"

# Summary
echo ""
echo "=== Results ==="
echo -e "  Passed: ${GREEN}$pass${NC}"
echo -e "  Failed: ${RED}$fail${NC}"

if [ "$fail" -gt 0 ]; then
  echo ""
  echo -e "${RED}FAILED:${NC}"
  for e in "${errors[@]}"; do
    echo "  - $e"
  done
  echo ""
  echo '{"success":false,"message":"Plugin package validation failed","code":"VERIFICATION_FAILED","details":{"passed":'"$pass"',"failed":'"$fail"'}}'
  exit 1
fi

echo ""
echo '{"success":true,"message":"Plugin package validation passed","code":"SUCCESS","details":{"passed":'"$pass"',"skills":'"$skill_count"'}}'
