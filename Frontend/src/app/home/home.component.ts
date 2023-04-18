import { Component, OnInit } from '@angular/core';
import {UserService} from "../services/user.service";

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  consola: any =
    {
      entrada: '',
      salida: ''
    }
  archivoSeleccionado: File|null = null;
  constructor(private service: UserService) { }
  ngOnInit(): void {
  }

  onFileSelected(event: Event) {
    const inputFile = event.target as HTMLInputElement;
    const archivo: File = inputFile.files![0];
    if (archivo && archivo.name.endsWith('.eea')) {
      this.archivoSeleccionado = archivo;
      const lector: FileReader = new FileReader();
      lector.readAsText(this.archivoSeleccionado);
      lector.onload = (e) => {
        if (e.target && e.target.result) {
          this.consola.entrada = e.target.result.toString()
        }
      }
    } else {
      alert('El archivo seleccionado no tiene extensión .eea');
    }
  }

  compilar() {
    console.log(this.consola)
    this.service.getIndex().subscribe(
      {
        next: (res) => {console.log(res)},
        error: (err) => {console.log(err)},
        complete: () => {}
      } )

    this.service.compilar(this.consola).subscribe(
      {
        next: (res: any) => {
          console.log(res);
          const JsonRespuesta = JSON.parse(JSON.stringify(res))
          this.consola.salida = JsonRespuesta.salida
        },
        error: (err: any) => { console.log(err) },
        complete: () => {}
      }
    )
  }
}
