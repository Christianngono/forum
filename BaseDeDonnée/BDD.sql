CREATE TABLE Utilisateur (
    ID INTEGER PRIMARY KEY,
    Nom_utilisateur TEXT,
    Email TEXT UNIQUE,
    Mot_de_passe TEXT
);

CREATE TABLE Publication (
    ID INTEGER PRIMARY KEY,
    Titre TEXT,
    Contenu TEXT,
    Utilisateur_ID INTEGER,
    Date_creation DATETIME,
    FOREIGN KEY(Utilisateur_ID) REFERENCES Utilisateur(ID)
);

CREATE TABLE Commentaire (
    ID INTEGER PRIMARY KEY,
    Contenu TEXT,
    Utilisateur_ID INTEGER,
    Publication_ID INTEGER,
    Date_creation DATETIME,
    FOREIGN KEY(Utilisateur_ID) REFERENCES Utilisateur(ID),
    FOREIGN KEY(Publication_ID) REFERENCES Publication(ID)
);

CREATE TABLE Categorie (
    ID INTEGER PRIMARY KEY,
    Nom TEXT
);
