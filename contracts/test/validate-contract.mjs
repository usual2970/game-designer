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
  const operationIds = [...content.matchAll(/operationId:\s+(\w+)/g)].map(m => m[1]);
  const schemaNames = [...content.matchAll(/^\s{4}(\w+):\n\s+type:\s+object/gm)].map(m => m[1]);
  const errorCodes = [...content.matchAll(/-\s+([A-Z_]+)/g)].map(m => m[1]);

  return { operationIds, schemaNames: [...new Set(schemaNames)], errorCodes: [...new Set(errorCodes)] };
}

const parsed = parseOpenAPI(contract);

// Test: contract has required info
assert(contract.includes('openapi: 3.0.3'), 'OpenAPI version is 3.0.3');
assert(contract.includes('title: Game Designer Server API'), 'Contract has title');
assert(contract.includes('slot-machine'), 'Contract describes slot-machine gameplay');
assert(!contract.includes('activity-style'), 'Contract does not describe generic activity-style games');

// Test: slot machine operations present
const requiredOperations = [
  'createOrResumeSession',
  'getPlayerProfile',
  'updatePlayerProfile',
  'getSlotConfig',
  'getBalance',
  'spin',
  'getSpinHistory',
  'getSlotLeaderboard',
];

for (const op of requiredOperations) {
  assert(parsed.operationIds.includes(op), `Operation ${op} is defined`);
}

// Test: legacy activity-game operations removed
const legacyOperations = ['getGameState', 'saveGameState', 'submitScore', 'getLeaderboard'];
for (const op of legacyOperations) {
  assert(!parsed.operationIds.includes(op), `Legacy operation ${op} is removed`);
}

// Test: all required schemas present
const requiredSchemas = [
  'SessionRequest',
  'SessionResponse',
  'ProfileResponse',
  'UpdateProfileRequest',
  'SlotConfigResponse',
  'BalanceResponse',
  'SpinRequest',
  'SpinResult',
  'ReelWindow',
  'PaylineWin',
  'SpinHistoryEntry',
  'SpinHistoryResponse',
  'SlotLeaderboardEntry',
  'SlotLeaderboardResponse',
  'SymbolDefinition',
  'PaylineDefinition',
  'Error',
];

for (const schema of requiredSchemas) {
  assert(
    contract.includes(`  ${schema}:`) || contract.includes(`${schema}:`),
    `Schema ${schema} is defined`
  );
}

// Test: legacy schemas removed
const legacySchemas = ['SaveGameStateRequest', 'GameStateResponse', 'SubmitScoreRequest', 'SubmitScoreResponse', 'LeaderboardEntry', 'LeaderboardResponse'];
for (const schema of legacySchemas) {
  assert(!contract.includes(`  ${schema}:`), `Legacy schema ${schema} is removed`);
}

// Test: error codes present including slot-specific ones
const requiredErrorCodes = [
  'INVALID_PARAMETERS',
  'UNAUTHORIZED',
  'NOT_FOUND',
  'SESSION_EXPIRED',
  'INSUFFICIENT_BALANCE',
  'INTERNAL_ERROR',
];

for (const code of requiredErrorCodes) {
  assert(contract.includes(code), `Error code ${code} is defined`);
}

// Test: session token parameter
assert(contract.includes('X-Session-Token'), 'Session token parameter is defined');

// Test: slot-specific tags
const requiredTags = ['Session', 'Profile', 'Slot', 'Balance', 'SpinHistory', 'Leaderboard'];
for (const tag of requiredTags) {
  assert(contract.includes(`tags: [${tag}]`) || contract.includes(`- name: ${tag}`), `Tag ${tag} is defined`);
}

// Test: legacy tags removed
const legacyTags = ['GameState', 'Score'];
for (const tag of legacyTags) {
  assert(!contract.includes(`tags: [${tag}]`) && !contract.includes(`- name: ${tag}`), `Legacy tag ${tag} is removed`);
}

// Test: spin result has required fields for server-authoritative resolution
assert(contract.includes('totalPayout'), 'Spin result includes totalPayout');
assert(contract.includes('paylineWins'), 'Spin result includes paylineWins');
assert(contract.includes('server-authoritative'), 'Contract specifies server-authoritative spin resolution');

// Summary
console.log(`\n${passCount} passed, ${failCount} failed`);
if (failCount > 0) {
  process.exit(1);
}
