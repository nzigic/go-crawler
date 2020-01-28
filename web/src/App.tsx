import React, { useState } from 'react';
import Button from '@material-ui/core/Button';
import './App.css';
import { makeStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';
import { CrawlerServiceClient } from './api/crawler-service';
import axios, { AxiosResponse } from "axios";
import { CrawlResult } from './api/crawler-vo';
import CircularProgress from '@material-ui/core/CircularProgress';
import { red } from '@material-ui/core/colors';

const useStyles = makeStyles({
  table: {
    minWidth: 650,
  },
  brokenLink: {
    borderLeft: '20px solid red'
  }
});

interface Row {
  url: string;
  broken: boolean;
  message?: string;
}

const defaultRows: Row[] = [];

const getTransport = (endpoint: string) => async function <T>(method: string, args: any[] | undefined) {
  return new Promise<T>(async (resolve, reject) => {
    try {
      let axiosPromise: AxiosResponse<T> = await axios.post<T>(
        endpoint + "/" + encodeURIComponent(method),
        JSON.stringify(args),
      );
      return resolve(axiosPromise.data);
    } catch (e) {
      return reject(e);
    }
  });
};

const fetchData = async () => {
  const svc = new CrawlerServiceClient(getTransport('http://localhost:8080' + CrawlerServiceClient.defaultEndpoint))
  const links = await svc.crawl('http://bestbytes.de');

  return links.map((link: CrawlResult): Row => ({
    url: link.Url,
    broken: link.Broken,
    message: link.Message
  }));
}

const App: React.FC = () => {
  const classes = useStyles();
  const [rows, setRows] = useState(defaultRows);
  const [isLoading, setIsLoading] = useState(false)

  const handleClick = async () => {
    setIsLoading(true);
    const rows = await fetchData();
    setRows(rows);
    setIsLoading(false)
  }

  return (
    <div className="App">
      <Button variant="contained" color="primary" onClick={handleClick}>
        Check bestbytes.de
        {isLoading && <CircularProgress color="secondary" />}
      </Button>

      <TableContainer component={Paper}>
        <Table className={classes.table} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableCell>Link</TableCell>
              <TableCell align="right">Link valid</TableCell>
              <TableCell align="right">Reason</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.map(row => (
              <TableRow key={row.url} className={row.broken ? classes.brokenLink : ''}>
                <TableCell component="th" scope="row">
                  {row.url}
                </TableCell>
                <TableCell align="right">{!row.broken ? 'Yes' : 'No'}</TableCell>
                <TableCell align="right">{row.message || '-'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
}

export default App;
