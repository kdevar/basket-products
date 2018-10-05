import React from 'react';
import {withStyles} from '@material-ui/core/styles';
import Geocode from "react-geocode";
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Divider from '@material-ui/core/Divider';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import TextField from '@material-ui/core/TextField';
import CircularProgress from '@material-ui/core/CircularProgress';

Geocode.setApiKey("AIzaSyByRxUnUyA9RKpU2FOuOtWkESHgf693xdo");

Geocode.enableDebug();

const geoCode = (val) => {
    return Geocode.fromAddress(val).then(
        response => {
            const {lat, lng} = response.results[0].geometry.location;
            return {lat, lng}
        },
        error => {
            console.error(error);
        }
    );
}


const styles = {
    card: {
        minWidth: 275,
    },
    bullet: {
        display: 'inline-block',
        margin: '0 2px',
        transform: 'scale(0.8)',
    },
    title: {
        marginBottom: 16,
        fontSize: 14,
    },
    pos: {
        marginBottom: 12,
    },
    progress: {
        margin: 50,
    }
};

class Location extends React.Component {
    constructor(props) {
        super(props);
        this.props = props;
        this.state = {
            locationText: ''
        }

    }

    componentWillReceiveProps(prevprops, nextprops) {

    }

    handleLocationChange() {
        return geoCode(this.state.locationText).then(loc => {
            return this.props.onLocationChange(loc)
        })
    }

    handleLocationTextChange(e) {
        return this.setState({locationText: e.target.value})
    }

    render() {
        const {classes} = this.props;
        const bull = <span className={classes.bullet}>•</span>;

        if (this.props.loading) {
            return (<CircularProgress className={classes.progress} size={50}/>)
        }

        return (
            <Card className={classes.card}>
                <CardHeader
                    action={
                        <IconButton>
                            <MoreVertIcon/>
                        </IconButton>
                    }
                    title="Shrimp and Chorizo Paella"
                    subheader="September 14, 2016"
                />

                <CardContent>
                    <form onSubmit={(event) => {
                        event.preventDefault();
                        this.handleLocationChange()
                    }}>
                        <TextField
                            id="standard-with-placeholder"
                            label="Switch to different location"
                            fullWidth={true}
                            placeholder="Switch to different location"
                            className={classes.textField}
                            margin="normal"
                            onChange={this.handleLocationTextChange.bind(this)}
                        />
                    </form>
                    <Typography className={classes.title} color="textSecondary">
                        Searching relative to
                    </Typography>
                    <Typography variant="subheading" component="h2">
                        Metro
                    </Typography>
                    <Typography className={classes.title} color="textSecondary">
                        {this.props.area.metroAreaName} ({this.props.area.metroAreaId})
                    </Typography>
                    <Typography variant="subheading" component="h2">
                        City
                    </Typography>
                    <Typography className={classes.title} color="textSecondary">
                        {this.props.area.cityName} ({this.props.area.cityId})
                    </Typography>
                    <Typography variant="subheading" component="h2">
                        Postal Code
                    </Typography>
                    <Typography className={classes.title} color="textSecondary">
                        {this.props.area.postalCode}
                    </Typography>

                    <Typography variant="subheading" component="h2">
                        Stores
                    </Typography>
                    <List>
                        {this.props.area.stores && this.props.area.stores.map(store => (
                            <div>
                                <ListItem disableGutters={true}>
                                    <ListItemText primary={store.FullStoreName} secondary={store.ChainDesc}/>
                                </ListItem>
                                <Divider/>
                            </div>

                        ))}
                    </List>

                </CardContent>

            </Card>
        );
    }
}


export default withStyles(styles)(Location);