'use strict';

const path = require('path');
const fs = require('fs');

function readFirstSentanceAndWriteItToResultFile(filename, pathToSourceDir, pathToDestinationDir) {
  return new Promise((resolve, reject) => {
    const pathToFile = `${pathToSourceDir}/${filename}`;
    // Create read stream and read file byte by byte
    const readStream = fs.createReadStream(pathToFile, { highWaterMark: 1 });

    let result = '';

    readStream.on('data', chunk => {
      result += chunk;

      // Regex for sentance separator
      const regexp = /([!?.]\s+)|(\.{3}\s+)/;

      const readFirstSentance = regexp.test(result);
      if (readFirstSentance) {
        readStream.close();

        const pathToOutputFile = `${pathToDestinationDir}/${filename}.res`;

        fs.writeFile(pathToOutputFile, result, (err) => {
          if (err) reject(err);
          resolve();
        });
      }
    });
  });
}

// Generate promises
function* generator(filenames, pathToSourceDir, pathToDestinationDir) {
  for (const filename of filenames) {
    yield readFirstSentanceAndWriteItToResultFile(
      filename, 
      pathToSourceDir, 
      pathToDestinationDir
    );
  }
}

(async function main() {
  let pathToSourceDir;
  let pathToDestinationDir;

  // Try to make an absolute paths from relative
  try {
    pathToSourceDir = path.resolve(process.argv[2]);
  } catch (err) {
    console.error('Wrong source directory path!');
    process.exit(1);
  }

  try {
    pathToDestinationDir = path.resolve(process.argv[3]);
  } catch (err) {
    console.error('Wrong destination directory path!');
    process.exit(1);
  }

  const sourceDirPathExists = fs.existsSync(pathToSourceDir);
  if (!sourceDirPathExists) {
    console.error('Source directory path does not exist!');
    process.exit(1);
  }

  const destinationDirPathExists = fs.existsSync(pathToDestinationDir);
  // Create destination dir if it does not exist
  if (!destinationDirPathExists) {
    try {
      fs.mkdirSync(pathToDestinationDir, { recursive: true });
    } catch (err) {
      console.error(err);
      process.exit(1);
    }
  }

  // Array of input file names
  const inputFiles = fs.readdirSync(pathToSourceDir);
  // Run in parallel 
  Promise.all(
    generator(inputFiles, pathToSourceDir, pathToDestinationDir)
  ).then(result => console.log(`Total number of processed files: ${result.length}`));
})();

