// src/features/dashboard/MetricsPanel.tsx
import { useQuery } from '@tanstack/react-query';
import {  Typography } from '@mui/material';
import axios from 'axios';

const fetchMetrics = async () => {
    // This endpoint would be a proxy to your Prometheus metrics
    const response = await axios.get('/api/metrics');
    return response.data;
};

const MetricsPanel = () => {
    const { data, isLoading } = useQuery({
        queryKey: ['metrics'],
        queryFn: fetchMetrics,
        refetchInterval: 30000, // Refresh every 30 seconds
    });

    if (isLoading) return <Typography>Loading metrics...</Typography>;

    return (
        <div>Metric Panel</div>
        // <Grid container spacing={3}>
        //     <Grid item xs={12} md={4}>
        //         <Card>
        //             <CardContent>
        //                 <Typography color="textSecondary" gutterBottom>
        //                     API Requests
        //                 </Typography>
        //                 <Typography variant="h4">
        //                     {data?.requestTotal || '0'}
        //                 </Typography>
        //             </CardContent>
        //         </Card>
        //     </Grid>
        //
        //     <Grid item xs={12} md={4}>
        //         <Card>
        //             <CardContent>
        //                 <Typography color="textSecondary" gutterBottom>
        //                     Average Response Time
        //                 </Typography>
        //                 <Typography variant="h4">
        //                     {data?.responseTime?.toFixed(2) || '0'} ms
        //                 </Typography>
        //             </CardContent>
        //         </Card>
        //     </Grid>
        //
        //     <Grid item xs={12} md={4}>
        //         <Card>
        //             <CardContent>
        //                 <Typography color="textSecondary" gutterBottom>
        //                     Error Rate
        //                 </Typography>
        //                 <Typography variant="h4">
        //                     {data?.errorRate?.toFixed(2) || '0'}%
        //                 </Typography>
        //             </CardContent>
        //         </Card>
        //     </Grid>
        // </Grid>
    );
};

export default MetricsPanel;