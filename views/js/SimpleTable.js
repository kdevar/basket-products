import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';

const styles = {
    root: {
        width: '100%',
        overflowX: 'auto'
    },
    table: {
        minWidth: 700,

    },
};

let id = 0;

function render(min,max){
    if(!min && !max) return "-";
    if (min === max) {
        return max;
    }
    return `${min}-${max}`;
}

function renderUserPrices(props){
    const { classes } = props;
    if(!props.data || props.data.length === 0){
        return (<div><Typography variant={"caption"}>No User or Cluster Prices available</Typography><Divider/></div>);
    }
    return (
    <Table className={classes.table}>
        <TableHead>
            <TableRow>
                <TableCell >Store Name</TableCell>
                <TableCell numeric>List Price</TableCell>
                <TableCell numeric>Sale Price</TableCell>
                <TableCell numeric>Final Price</TableCell>
            </TableRow>
        </TableHead>
        <TableBody>
            {props.data && props.data.map(price => (
                <TableRow>
                    <TableCell >{price.fullStoreName}</TableCell>

                    <TableCell numeric>{price.listPrice}</TableCell>
                    <TableCell numeric>{price.salePrice}</TableCell>
                    <TableCell numeric>{price.finalPrice}</TableCell>
                </TableRow>
            ))}
        </TableBody>
    </Table>);
}

function SimpleTable(props) {
    const { classes } = props;

    return (
        <div className={classes.root}>
            <Typography variant={"subheading"}>User or Cluster based Price</Typography>

            {renderUserPrices(props)}

            <br/>

            <Typography variant={"subheading"}>Estimated Prices</Typography>

            <Table className={classes.table}>
                <TableHead>
                    <TableRow>
                        <TableCell >Store Name</TableCell>

                        <TableCell numeric>City Est</TableCell>
                        <TableCell numeric>Metro Est</TableCell>
                        <TableCell numeric>Zip Est</TableCell>
                        <TableCell numeric>0-50mi Est</TableCell>
                        <TableCell numeric>50-100mi Est</TableCell>
                        <TableCell numeric>100-4000mi Est</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {props.estimatedData && props.estimatedData.map(price => (
                        <TableRow>
                            <TableCell >{price.FullStoreName}</TableCell>
                            <TableCell numeric>{render(price.CITY.Min,price.CITY.Max)}</TableCell>
                            <TableCell numeric>{render(price.METRO.Min,price.METRO.Max)}</TableCell>
                            <TableCell numeric>{render(price.ZIP.Min,price.ZIP.Max)}</TableCell>
                            <TableCell numeric>{render(price.FIFTYMILE.Min,price.FIFTYMILE.Max)}</TableCell>
                            <TableCell numeric>{render(price.HUNDREDMILES.Min,price.HUNDREDMILES.Max)}</TableCell>
                            <TableCell numeric>{render(price.NATIONALMILES.Min,price.NATIONALMILES.Max)}</TableCell>

                        </TableRow>
                    ))}
                </TableBody>
            </Table>


        </div>
    );
}

SimpleTable.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(SimpleTable);