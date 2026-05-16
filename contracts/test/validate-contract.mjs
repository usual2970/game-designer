import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const contractPath = join(__dirname, '..', 'game-server.openapi.yaml');

const contract = readFileSync(contractPath, 'utf-8');

let passCount = 0;
let failCount = 0;

function assert(condition, message) {
  if (condition) {
    passCount++;
  } else {
    failCount++;
    console.error(`FAIL: ${message}`);
  }
}

// Parse the YAML contract (simple regex-based validation for MVP)
function parseOpenAPI(content) {
  // Extract paths section
  const pathsMatch = content.match(/^paths:\n((?:  .+\n?)*)/m);
  const paths = pathsMatch ? pathsMatch[1] : '';

  // Extract operation IDs
  const operationIds = [...content.matchAll(/operationId:\s+(\w+)/g)].map(m => m[1]);

  // Extract schemas
  const schemaNames = [...content.matchAll(/^\s{4}(\w+):\n\s+type:\s+object/gm)].map(m => m[1]);

  // Extract error codes
  const errorCodes = [...content.matchAll(/-\s+([A-Z_]+)/g)].map(m => m[1]);

  return { paths, operationIds, schemaNames: [...new Set(schemaNames)], errorCodes: [...new Set(errorCodes)] };
}

const parsed = parseOpenAPI(contract);

// Test: contract has required info
assert(contract.includes('openapi: 3.0.3'), 'OpenAPI version is 3.0.3');
assert(contract.includes('title: Game Designer Server API'), 'Contract has title');

// Test: all MVP operations present
const requiredOperations = [
  'createOrResumeSession',
  'getPlayerProfile',
  'updatePlayerProfile',
  'getGameState',
  'saveGameState',
  'submitScore',
  'getLeaderboard',
];

for (const op of requiredOperations) {
  assert(parsed.operationIds.includes(op), `Operation ${op} is defined`);
}

// Test: all required schemas present
const requiredSchemas = [
  'SessionRequest',
  'SessionResponse',
  'ProfileResponse',
  'UpdateProfileRequest',
  'SaveGameStateRequest',
  'GameStateResponse',
  'SubmitScoreRequest',
  'SubmitScoreResponse',
  'LeaderboardEntry',
  'LeaderboardResponse',
  'Error',
];

for (const schema of requiredSchemas) {
  assert(
    contract.includes(`  ${schema}:`) || contract.includes(`${schema}:`),
    `Schema ${schema} is defined`
  );
}

// Test: error codes present
const requiredErrorCodes = [
  'INVALID_PARAMETERS',
  'UNAUTHORIZED',
  'NOT_FOUND',
  'SESSION_EXPIRED',
  'INTERNAL_ERROR',
];

for (const code of requiredErrorCodes) {
  assert(contract.includes(code), `Error code ${code} is defined`);
}

// Test: session token parameter
assert(contract.includes('X-Session-Token'), 'Session token parameter is defined');

// Test: tags cover all MVP capability groups
const requiredTags = ['Session', 'Profile', 'GameState', 'Score', 'Leaderboard'];
for (const tag of requiredTags) {
  assert(contract.includes(`tags: [${tag}]`) || contract.includes(`- name: ${tag}`), `Tag ${tag} is defined`);
}

// Summary
console.log(`\n${passCount} passed, ${failCount} failed`);
if (failCount > 0) {
  process.exit(1);
}
