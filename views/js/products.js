import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import GridListTileBar from '@material-ui/core/GridListTileBar';
import IconButton from '@material-ui/core/IconButton';
import StarBorderIcon from '@material-ui/icons/StarBorder';
import Typography from '@material-ui/core/Typography';
import classNames from 'classnames';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelActions from '@material-ui/core/ExpansionPanelActions';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Chip from '@material-ui/core/Chip';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import SimpleTable from './SimpleTable';

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
    icon: {
        verticalAlign: 'bottom',
        height: 20,
        width: 20,
    },
    details: {
        alignItems: 'center',
    },
    column: {
        flexBasis: '25%',
    },
    helper: {
        borderLeft: `2px solid ${theme.palette.divider}`,
        padding: `${theme.spacing.unit}px ${theme.spacing.unit * 2}px`,
    },
    link: {
        color: theme.palette.primary.main,
        textDecoration: 'none',
        '&:hover': {
            textDecoration: 'underline',
        },
    },
});

const Products = props => {
    const {classes} = props;
    return (
        <div className={classes.root}>
            {props.tileData &&
             props.tileData.map(tile => <StyledExpansionRow
                 getStores={props.getStores}
                 getStoresQuery={props.getStoresQuery}
                 getChainsQuery={props.getChainsQuery}
                 getLatitude={props.getLatitude}
                 getLongitude={props.getLongitude}
                 getLocationData={props.getLocationData}
                 getProductData={props.getProductData} row={tile} />)}
        </div>
    )
}

class ExpansionRow extends React.Component{
    constructor(props){
        super(props);
        this.props = props;
        this.state = {
            details: [],
            estimateDetails: null
        }
    }

    getProductData(productRow){
        return fetch(`/api/basket-products/${productRow.productId}/prices?${this.props.getStoresQuery()}`, {
            headers: {
                "latitude": this.props.getLatitude(),
                "longitude": this.props.getLongitude()
            },
        }).then(response => response.json())
    }

    getProductEstimateData(productRow){
        const location = this.props.getLocationData();

        return fetch(`/api/basket-products/${productRow.productId}/estimated-prices?${this.props.getChainsQuery()}&cityId=${location.cityId}&zipCodeId=${location.postalCodeId}&metroAreaId=${location.metroAreaId}`, {
            headers: {
                "latitude": this.props.getLatitude(),
                "longitude": this.props.getLongitude()
            },
        }).then(response => response.json())
    }

    onChange(event, expanded){
        if(expanded){
            this.getProductData(this.props.row).then(details => this.setState({details}))
            this.getProductEstimateData(this.props.row).then(estimateDetails => this.setState({estimateDetails}))
        }
    }
    
    render(){
        const {classes} = this.props;
        return (
            <ExpansionPanel onChange={this.onChange.bind(this)}>
                <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
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
                <StyleExpansionDetails estimateDetails={this.state.estimateDetails} getStores={this.props.getStores} data={this.state.details}/>
                <Divider />
            </ExpansionPanel>
        );
    }
}

const ExpansionDetails = props => {
    const {classes} = props;
    const stores = props.getStores();
    let estimatedData = [];
    if(props.estimateDetails){
         estimatedData = Object.keys(props.estimateDetails).map(chainId => {

            const store = stores &&!storeExists && stores.find(s => s.ChainID == chainId);
            const storeExists = store && props.data && props.data.find(s => s.storeId === store.StoreID)
            if(store && !storeExists){
                return {
                    ...store,
                    ...props.estimateDetails[chainId]
                }
            }
            return null;

        }).filter(n => n);
    }

    return (<ExpansionPanelDetails className={classes.details}>
        <SimpleTable estimatedData={estimatedData} data ={props.data}/>
    </ExpansionPanelDetails>)
};

const StyledExpansionRow = withStyles(styles)(ExpansionRow)
const StyleExpansionDetails = withStyles(styles)(ExpansionDetails)

export default withStyles(styles)(Products);