import { writeFileSync, readFileSync } from "fs";
import { unparse } from 'papaparse'

interface Output {
  [asciiCharCode: string]: number;
}

interface ResultEntry {
  character: string;
  charCode: number;
  occurrences: number;
}

const output: Output = JSON.parse(
  readFileSync("./output.json", { encoding: "utf-8" })
);

const outputAsTuple = Array.from(Object.entries(output));

let letterOccurrences = 0;
let numberOccurrences = 0;

const result: ResultEntry[] = [];

// For each charcode, push it to the results array, excepts if it's a digit (0 - 9) or a letter (A - Z or a - z), those are added at the end 
for (const [_char, occurrences] of outputAsTuple) {
  const char = Number(_char);

  // 0 - 9
  if (char >= 48 && char <= 57) {
    numberOccurrences += occurrences;
    continue;
  }

  if (
    (char >= 65 && char <= 90) || // A - Z
    (char >= 97 && char <= 122)   // a - z)
  ) {
    letterOccurrences += occurrences;
    continue;
  }

  result.push({
    character: String.fromCharCode(char),
    charCode: char,
    occurrences,
  });
}

result.push({
  character: "0 - 9",
  charCode: NaN,
  occurrences: numberOccurrences,
});

result.push({
  character: "A - z",
  charCode: NaN,
  occurrences: letterOccurrences,
});

// Sort by occurrences
const sortedResult = result.sort((a, b) =>
  a.occurrences > b.occurrences ? -1 : 1
);

console.log(sortedResult.map(x => ({ ...x, character: JSON.stringify(x.character) })))
console.table(sortedResult);
writeFileSync('./output.csv', unparse(sortedResult));