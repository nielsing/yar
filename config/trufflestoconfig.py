import json
import argparse

NOISELEVEL = 3

def read_truffles(filename):
    truffles = json.loads(filename.read())
    return truffles

def create_config(truffles):
    config = {'Rules': [], 'FileBlacklist': []}
    for key, value in truffles.items():
        config['Rules'].append({
            'Reason': key,
            'Rule': value,
            'Noise': NOISELEVEL
        })
    return config

def write_config(config):
    with open('yarconfig.json', 'w') as f:
        json.dump(config, f, indent=4)

def main(filename):
    truffles = read_truffles(filename)
    config = create_config(truffles)
    write_config(config)
    filename.close()

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Process a given truffleHog config file to a yar config file')
    parser.add_argument('filename', type=argparse.FileType('r'))
    args = parser.parse_args()
    main(args.filename)
