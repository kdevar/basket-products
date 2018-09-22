import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { withStyles } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Drawer from '@material-ui/core/Drawer';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import List from '@material-ui/core/List';
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';
import IconButton from '@material-ui/core/IconButton';
import Badge from '@material-ui/core/Badge';
import MenuIcon from '@material-ui/icons/Menu';
import ChevronLeftIcon from '@material-ui/icons/ChevronLeft';
import NotificationsIcon from '@material-ui/icons/Notifications';
import { mainListItems, secondaryListItems } from './listItems';
import SearchIcon from '@material-ui/icons/Search';
import { fade } from '@material-ui/core/styles/colorManipulator';
import SimpleTable from './SimpleTable';
import Typeahead from './typeahead';
import Input from '@material-ui/core/Input';
import Products from './products';
import TextField from '@material-ui/core/TextField'
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import Location from './location'



const drawerWidth = 240;

const styles = theme => ({
    root: {
        display: 'flex',
    },

    search:{
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
        width: '100%' ,
    },
},
    placeholder:{
        paddingLeft: theme.spacing.unit * 10,
        color: 'white'
    },
    toolbarIcon: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'flex-end',
        padding: '0 8px',
        ...theme.mixins.toolbar,
    },
    appBar: {
        zIndex: theme.zIndex.drawer + 1,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
    },
    appBarShift: {
        marginLeft: drawerWidth,
        width: `calc(100% - ${drawerWidth}px)`,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
        }),
    },
    menuButton: {
        marginLeft: 12,
        marginRight: 36,
    },
    menuButtonHidden: {
        display: 'none',
    },
    title: {
        flexGrow: 1,
    },
    drawerPaper: {
        position: 'relative',
        whiteSpace: 'nowrap',
        width: drawerWidth,
        transition: theme.transitions.create('width', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
        }),
    },
    drawerPaperClose: {
        overflowX: 'hidden',
        transition: theme.transitions.create('width', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
        width: theme.spacing.unit * 7,
        [theme.breakpoints.up('sm')]: {
            width: theme.spacing.unit * 9,
        },
    },
    appBarSpacer: theme.mixins.toolbar,
    content: {
        flexGrow: 1,
        padding: theme.spacing.unit * 3,
        paddingTop: theme.spacing.unit * 5,
        height: '100vh',
        overflow: 'auto',
    },
    chartContainer: {
        marginLeft: -22,
    },
    tableContainer: {
        height: 320,
    },
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
        response:[],
        location: [],
        loadingLocation: true,
        productIdDetails:{

        },
        productDetails: {

        }
    };

    onSearchCompleted(suggestion){
        this.setState({response: []});
        return this.fetchItems(suggestion);
    }

    componentDidMount(){
        this.fetchLocationDetails();
    }

    fetchLocationDetails(loc){
        this.setState({loadingLocation:true, response:[]});
        if(loc){
            this.setState({
                latitude: loc.lat,
                longitude: loc.lng
            });
        }
        fetch("/api/area", {
            headers:{
                latitude: this.state.latitude,
                longitude: this.state.longitude
            }
        }).
        then(response => response.json()).then(response => {
            this.setState({location:response, loadingLocation:false})
            return this.fetchItems();
        })
    }

    fetchItems(s){
        if(s){
            this.setState({suggestion:s})
        }
        const suggestion = s || this.state.suggestion;

        if(!suggestion || !suggestion.name){
            return;
        }

        return fetch(`/api/basket-products/?keyword=${suggestion.name}${suggestion.type === "Type" ? "&typeId="+suggestion.id : ""}${suggestion.type === "Brand" ? "&brandId="+suggestion.id : ""}&category=${suggestion.category || ""}`, {
            headers: {
                "latitude": this.state.latitude,
                "longitude": this.state.longitude
            },
        }).then(r => r.json()).then(response => {
            this.setState({response})
        })
    }
    getProductData(productRow){
        fetch(`/api/basket-products/${productRow.productId}/prices`, {
            headers: {
                "latitude": this.state.latitude,
                "longitude": this.state.longitude
            },
        }).then(response => {

        })
    }
    getProductEstimateData(productRow){
        fetch(`/api/basket-products/${productRow.productId}/?${getChainsAsUrlParams(this.state.location.chains)}&metroAreaId=${this.state.location.metroAreaId}&cityId=${this.state.location.cityId}&zipCodeId=${this.state.location.zipCodeId}`, {
            headers: {
                "latitude": this.state.latitude,
                "longitude": this.state.longitude
            },
        }).then(response => {

        })
    }

    handleDrawerOpen = () => {
        this.setState({ open: true });
    };

    handleDrawerClose = () => {
        this.setState({ open: false });
    };

    render() {
        const { classes } = this.props;

        return (
            <React.Fragment>
                <CssBaseline />
                <div className={classes.root}>
                    <AppBar
                        position="absolute"
                        className={classNames(classes.appBar, this.state.open && classes.appBarShift)}>
                        <Toolbar disableGutters={true}>
                            <Grid container spacing={24}>
                                <Grid item xs={12}>
                                    <div className={classes.search}>
                                        <div className={classes.searchIcon}>
                                            <SearchIcon />
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


                        <div className={classes.appBarSpacer} />

                        <Grid container spacing={24}>
                            <Grid item xs={9}>
                                <Products getStores={() => this.state.location.stores}
                                          getStoresQuery={() => getStoresAsUrlParams(this.state.location.stores)}
                                          getChainsQuery={() => getChainsAsUrlParams(this.state.location.chains)}
                                          getLatitude={() => this.state.latitude}
                                          getLongitude={() => this.state.longitude}
                                          getProductData={this.getProductData.bind(this)}
                                          getLocationData={() => this.state.location}
                                          tileData={this.state.response}/>
                            </Grid>
                            <Grid item xs={3}>
                                <Location loading={this.state.loadingLocation} area={this.state.location} onLocationChange={this.fetchLocationDetails.bind(this)}/>
                            </Grid>
                        </Grid>



                        <div className={classes.tableContainer}>

                        </div>
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