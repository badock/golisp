import { Component } from '@angular/core';
import {LispService} from "./lisp.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'webapp';

  lispCode = "";
  selectedExample = "";
  resultExpression = "";
  lispService: LispService;
  exampleNames = [];

  constructor(lispService: LispService) {
    this.lispService = lispService
    this.exampleNames = Object.keys(lispService.examples);

    if (this.exampleNames.length > 0) {
      this.selectedExample = this.exampleNames[0];
      this.selectedExampleHasChanged();
    }
  }

  runCode = () => {
    console.log("I will run the following expression:");
    console.log(this.lispCode);
    let lispCodeWithoutNewLines = this.lispCode.replace("\n", "");
    this.lispService.evalLispCode(lispCodeWithoutNewLines).subscribe((result) => {
      // @ts-ignore
      this.resultExpression = `${result.result}`
    })
  }

  selectedExampleHasChanged = () => {
    console.log(`selectedExample => ${this.selectedExample}`);
    this.lispCode = this.lispService.examples[this.selectedExample];
  }

}
