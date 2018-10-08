module.exports = {
    entry: './index.js',
    devtool: 'sourcemap',
    watch: false,
    output: {
        path: __dirname + '/dist',
        publicPath: '/',
        filename: 'app.min.js'
    },
    module: {
        rules: [
            {
                test: /\.js$/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader"
                }
            }
        ]
    }
};