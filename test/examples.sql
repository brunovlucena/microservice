CREATE TABLE configs (
  id SERIAL,
  data JSONB
);

-- Make sure Name is unique
CREATE UNIQUE INDEX configs_name_idx ON configs(((data->>'name')::name));

-- Example
CREATE INDEX idxmon ON configs ((data->>'monitoring'));

INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-0",
    "metadata": {
      "monitoring": {
        "enabled": false
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "128m"
        },
        "memory": {
          "enabled": true,
          "value": "1024Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-1",
    "metadata": {
      "monitoring": {
        "enabled": true
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "128m"
        },
        "memory": {
          "enabled": true,
          "value": "512Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-2",
    "metadata": {
      "monitoring": {
        "enabled": false
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "512m"
        },
        "memory": {
          "enabled": false,
          "value": "1024Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-3",
    "metadata": {
      "monitoring": {
        "enabled": true
      },
      "limits": {
        "cpu": {
          "enabled": false,
          "value": "512m"
        },
        "memory": {
          "enabled": true,
          "value": "2048Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-4",
    "metadata": {
      "monitoring": {
        "enabled": true
      },
      "limits": {
        "cpu": {
          "enabled": false,
          "value": "1024m"
        },
        "memory": {
          "enabled": true,
          "value": "2048Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-5",
    "metadata": {
      "monitoring": {
        "enabled": true
      },
      "limits": {
        "cpu": {
          "enabled": false,
          "value": "256m"
        },
        "memory": {
          "enabled": true,
          "value": "256Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-6",
    "metadata": {
      "monitoring": {
        "enabled": false
      },
      "limits": {
        "cpu": {
          "enabled": false,
          "value": "512m"
        },
        "memory": {
          "enabled": false,
          "value": "512Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-7",
    "metadata": {
      "monitoring": {
        "enabled": true
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "256m"
        },
        "memory": {
          "enabled": false,
          "value": "256Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-8",
    "metadata": {
      "monitoring": {
        "enabled": false
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "256m"
        },
        "memory": {
          "enabled": true,
          "value": "1024Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-9",
    "metadata": {
      "monitoring": {
        "enabled": true
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "1024m"
        },
        "memory": {
          "enabled": false,
          "value": "512Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-10",
    "metadata": {
      "monitoring": {
        "enabled": true
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "1024m"
        },
        "memory": {
          "enabled": true,
          "value": "256Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-11",
    "metadata": {
      "monitoring": {
        "enabled": true
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "256m"
        },
        "memory": {
          "enabled": true,
          "value": "512Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-12",
    "metadata": {
      "monitoring": {
        "enabled": true
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "512m"
        },
        "memory": {
          "enabled": true,
          "value": "512Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-13",
    "metadata": {
      "monitoring": {
        "enabled": false
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "256m"
        },
        "memory": {
          "enabled": true,
          "value": "512Mi"
        }
      }
    }
  }');
INSERT INTO configs (data) VALUES ('
  {
    "name": "pod-14",
    "metadata": {
      "monitoring": {
        "enabled": false
      },
      "limits": {
        "cpu": {
          "enabled": true,
          "value": "512m"
        },
        "memory": {
          "enabled": false,
          "value": "1024Mi"
        }
      }
    }
  }');
