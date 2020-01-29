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
import { CrawlResult } from './api/crawler-vo';
import getClient from './utils/getClient';
import { LinearProgress } from '@material-ui/core';

const useStyles = makeStyles({
  table: {
    minWidth: 650,
  },
  brokenLink: {
    borderLeft: '20px solid red'
  }
});

interface Row {
  url?: string;
  broken?: boolean;
  message?: string;
}

const defaultRows: Row[] = [];

const fetchData = async () => {
  const client = getClient();
  const links = await client.crawl('http://bestbytes.de');

  return links.map((link: CrawlResult): Row => ({
    url: link.url,
    broken: link.broken,
    message: link.message
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
      </Button>
      {isLoading && <LinearProgress color="primary" />}

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
                  <a href={row.url} target="_blank">{row.url}</a>
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
