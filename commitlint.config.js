/* eslint-disable */
module.exports = {
    rules: {
        'header-max-length': [2, 'always', 100],
        'body-max-line-length': [0, 'always', 100],
        'type-enum': [
            2,
            'always',
            [
                'build',
                'chore',
                'ci',
                'docs',
                'feat',
                'fix',
                'perf',
                'refactor',
                'revert',
                'style',
                'test',
                'debug'
            ],
        ],
    },
    extends: ['@commitlint/config-conventional']
};
