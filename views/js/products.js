import React from 'react';
import {withStyles} from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Divider from '@material-ui/core/Divider';
import ProductDetails from './product-details';

const styles = theme => ({
    root: {
        width: '100%',
    },
    heading: {
        fontSize: theme.typography.pxToRem(15),
    },
    secondaryHeading: {
        fontSize: theme.typography.pxToRem(15),
        color: theme.palette.text.secondary,
    },
    details: {
        alignItems: 'center',
    },
    column: {
        flexBasis: '25%',
    }
});

const storeList = stores => stores.map(s => s.StoreID).join(",");
const chainList = chains => {
    return Object.keys(chains).map(key => chains[key].ChainID).join(",");
}

const Products = props => {
    const {classes} = props;
    return (
        <div className={classes.root}>
            {props.rowData && props.rowData.map(tile => <StyledExpansionRow
                getStores={props.getStores}
                getChains={props.getChains}
                getLatitude={props.getLatitude}
                getLongitude={props.getLongitude}
                getLocationData={props.getLocationData}
                row={tile}/>)
            }
        </div>
    )
}

class ExpansionRow extends React.Component {
    constructor(props) {
        super(props);
        this.props = props;
        this.state = {
            details: [],
            estimateDetails: null
        }
    }

    getProductData(productRow) {
        const stores = storeList(this.props.getStores());
        return fetch(`/api/basket-products/${productRow.productId}/prices?storeId=${stores}`, {
            headers: {
                "latitude": this.props.getLatitude(),
                "longitude": this.props.getLongitude()
            },
        }).then(response => response.json())
    }

    getProductEstimateData(productRow) {
        const location = this.props.getLocationData();
        const chains = chainList(this.props.getChains());
        return fetch(`/api/basket-products/${productRow.productId}/estimated-prices?chainId=${chains}&cityId=${location.cityId}&zipCodeId=${location.postalCodeId}&metroAreaId=${location.metroAreaId}`, {
            headers: {
                "latitude": this.props.getLatitude(),
                "longitude": this.props.getLongitude()
            },
        }).then(response => response.json())
    }

    onChange(event, expanded) {
        if (expanded) {
            this.getProductData(this.props.row)
                .then(details => this.setState({details}))
            this.getProductEstimateData(this.props.row)
                .then(estimateDetails => this.setState({estimateDetails}))
        }
    }

    render() {
        const {classes} = this.props;
        return (
            <ExpansionPanel onChange={this.onChange.bind(this)}>
                <ExpansionPanelSummary expandIcon={<ExpandMoreIcon/>}>
                    <div className={classes.column}>
                        <Typography className={classes.secondaryHeading}>{this.props.row.productId}</Typography>
                    </div>
                    <div className={classes.column}>
                        <Typography className={classes.secondaryHeading}>{this.props.row.typeDesc}</Typography>
                    </div>
                    <div className={classes.column}>
                        <Typography className={classes.heading}>{this.props.row.brandDesc}</Typography>
                    </div>
                    <div className={classes.column}>
                        <Typography className={classes.secondaryHeading}>{this.props.row.productDesc}</Typography>
                    </div>
                    <div className={classes.column}>
                        <Typography className={classes.secondaryHeading}>{this.props.row.sizeDesc}</Typography>
                    </div>
                </ExpansionPanelSummary>
                <StyleExpansionDetails estimateDetails={this.state.estimateDetails} getStores={this.props.getStores}
                                       data={this.state.details}/>
                <Divider/>
            </ExpansionPanel>
        );
    }
}

const ExpansionDetails = props => {
    const {classes} = props;
    const stores = props.getStores();
    let estimatedData = [];
    if (props.estimateDetails) {
        estimatedData = Object
            .keys(props.estimateDetails)
            .map(chainId => {
                const store = stores && stores.find(s => s.ChainID == chainId);
                const storeExistsInLiveData = store && props.data && props.data.find(s => s.storeId === store.StoreID)

                if (store && !storeExistsInLiveData) {
                    return {
                        ...store,
                        ...props.estimateDetails[chainId]
                    }
                }

                return null;

            }).filter(n => n);
    }

    return (
        <ExpansionPanelDetails className={classes.details}>
            <ProductDetails
                estimatedData={estimatedData}
                liveData={props.data}/>
        </ExpansionPanelDetails>
    )
};

const StyledExpansionRow = withStyles(styles)(ExpansionRow)
const StyleExpansionDetails = withStyles(styles)(ExpansionDetails)

export default withStyles(styles)(Products);