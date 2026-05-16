#!/usr/bin/env bash
# verify-plugin-package.sh — Validate the Game Designer plugin package structure
#
# Usage: ./scripts/verify-plugin-package.sh
#
# Checks manifests, skills, bundled assets, and documentation references.
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
SKILLS_DIR="$ROOT_DIR/plugin/skills"
check "Skills directory exists" test -d "$SKILLS_DIR"

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
check "plugin/README.md exists" test -f "$ROOT_DIR/plugin/README.md"
check "plugin/INSTALL.md exists" test -f "$ROOT_DIR/plugin/INSTALL.md"
check "docs/integration/plugin-installation.md exists" test -f "$ROOT_DIR/docs/integration/plugin-installation.md"
check "docs/integration/agent-golden-path.md exists" test -f "$ROOT_DIR/docs/integration/agent-golden-path.md"
check "cli/README.md references setup skill" grep -q "setup-game-designer-cli" "$ROOT_DIR/cli/README.md"

# 7. Install docs do not claim CLI is built during install
echo ""
echo "7. Install docs accuracy"
check "INSTALL.md does not claim auto-CLI-build" python3 -c "
text = open('$ROOT_DIR/plugin/INSTALL.md').read()
assert 'does not compile' in text.lower() or 'does not build' in text.lower()
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
