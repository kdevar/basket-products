import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import {withStyles} from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import SearchIcon from '@material-ui/icons/Search';
import Typeahead from './typeahead';
import Products from './products';
import Grid from '@material-ui/core/Grid';
import Location from './location';
import { fade } from '@material-ui/core/styles/colorManipulator';

const styles = theme => ({
    root: {
        display: 'flex',
    },

    search: {
        position: 'relative',
        borderRadius: theme.shape.borderRadius,
        backgroundColor: fade(theme.palette.common.white, 0.15),
        '&:hover': {
            backgroundColor: fade(theme.palette.common.white, 0.25),
        },
        marginRight: theme.spacing.unit * 2,
        marginLeft: 0,
        width: '100%',
        [theme.breakpoints.up('sm')]: {
            marginLeft: theme.spacing.unit * 3,
            width: 'auto',
        },
    },

    searchIcon: {
        width: theme.spacing.unit * 9,
        height: '100%',
        position: 'absolute',
        pointerEvents: 'none',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
    inputRoot: {
        color: 'inherit',
        width: '100%',
    },
    inputInput: {
        paddingTop: theme.spacing.unit,
        paddingRight: theme.spacing.unit,
        paddingBottom: theme.spacing.unit,
        paddingLeft: theme.spacing.unit * 10,
        transition: theme.transitions.create('width'),
        width: '100%',
        [theme.breakpoints.up('md')]: {
            width: '100%',
        },
    },
    placeholder: {
        paddingLeft: theme.spacing.unit * 10,
        color: 'white'
    },
    appBar: {
        zIndex: theme.zIndex.drawer + 1,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
    },
    appBarSpacer: theme.mixins.toolbar,
    content: {
        flexGrow: 1,
        padding: theme.spacing.unit * 3,
        paddingTop: theme.spacing.unit * 5,
        height: '100vh',
        overflow: 'auto',
    }
});

const getChainsAsUrlParams = (chains) => {
    return Object.keys(chains).map(key => `chainId=${chains[key].ChainID}`).join("&")
}

const getStoresAsUrlParams = (stores) => {
    return stores.map(store => `storeId=${store.StoreID}`).join("&")
}

class Dashboard extends React.Component {
    state = {
        open: false,
        latitude: "38.8876531",
        longitude: "-77.0954574",
        response: [],
        location: [],
        loadingLocation: true,
        productIdDetails: {},
        productDetails: {}
    };

    onSearchCompleted(suggestion) {
        this.setState({response: []});
        return this.fetchItems(suggestion);
    }

    componentDidMount() {
        this.fetchLocationDetails();
    }

    fetchLocationDetails(loc) {
        this.setState({loadingLocation: true, response: []});
        if (loc) {
            this.setState({
                latitude: loc.lat,
                longitude: loc.lng
            });
        }
        fetch("/api/area", {
            headers: {
                latitude: this.state.latitude,
                longitude: this.state.longitude
            }
        }).then(response => response.json()).then(response => {
            this.setState({location: response, loadingLocation: false})
            return this.fetchItems();
        })
    }

    fetchItems(s) {
        if (s) {
            this.setState({suggestion: s})
        }
        const suggestion = s || this.state.suggestion;

        if (!suggestion || !suggestion.name) {
            return;
        }

        return fetch(`/api/basket-products/?keyword=${suggestion.name}${suggestion.type === "Type" ? "&typeId=" + suggestion.id : ""}${suggestion.type === "Brand" ? "&brandId=" + suggestion.id : ""}&category=${suggestion.category || ""}`, {
            headers: {
                "latitude": this.state.latitude,
                "longitude": this.state.longitude
            },
        }).then(r => r.json()).then(response => {
            this.setState({response})
        })
    }

    render() {
        const {classes} = this.props;

        return (
            <React.Fragment>
                <CssBaseline/>
                <div className={classes.root}>
                    <AppBar
                        position="absolute"
                        className={classes.AppBar}>
                        <Toolbar disableGutters={true}>
                            <Grid container spacing={24}>
                                <Grid item xs={12}>
                                    <div className={classes.search}>
                                        <div className={classes.searchIcon}>
                                            <SearchIcon/>
                                        </div>

                                        <Typeahead
                                            latitude={this.state.latitude}
                                            longitude={this.state.longitude}
                                            onSelected={this.onSearchCompleted.bind(this)}
                                            classes={{
                                                root: classes.inputRoot,
                                                input: classes.inputInput,
                                                placeholder: classes.placeholder
                                            }}
                                        />
                                    </div>
                                </Grid>
                            </Grid>

                        </Toolbar>
                    </AppBar>

                    <main className={classes.content}>

                        <div className={classes.appBarSpacer}/>

                        <Grid container spacing={24}>
                            <Grid item xs={9}>
                                <Products getStores={() => this.state.location.stores}
                                          getStoresQuery={() => getStoresAsUrlParams(this.state.location.stores)}
                                          getChainsQuery={() => getChainsAsUrlParams(this.state.location.chains)}
                                          getLatitude={() => this.state.latitude}
                                          getLongitude={() => this.state.longitude}
                                          getLocationData={() => this.state.location}
                                          tileData={this.state.response}/>
                            </Grid>
                            <Grid item xs={3}>
                                <Location loading={this.state.loadingLocation} area={this.state.location}
                                          onLocationChange={this.fetchLocationDetails.bind(this)}/>
                            </Grid>
                        </Grid>
                    </main>
                </div>
            </React.Fragment>
        );
    }
}

Dashboard.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Dashboard);