import React from 'react';
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';
import { Link } from 'react-router-dom';

export default function App() {
  return (
    <Navbar bg="light" expand="lg">
      <Navbar.Brand href="/">Bevy</Navbar.Brand>
      <Navbar.Toggle aria-controls="basic-navbar-nav" />
      <Navbar.Collapse id="basic-navbar-nav">
        <Nav className="mr-auto">
          <Link
            to={{
              pathname: '/movies',
            }}
            className="nav-link"
          >
            Movies
          </Link>
          <Link
            to={{
              pathname: '/torrents',
            }}
            className="nav-link"
          >
            Torrents
          </Link>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  );
}
